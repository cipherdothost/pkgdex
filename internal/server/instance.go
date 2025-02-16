// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"time"

	"git.sr.ht/~jamesponddotco/xstd-go/xcrypto/xtls"
	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp/xmiddleware"
	"go.cipher.host/cmdkit"
	"go.cipher.host/pkgdex/internal/config"
	"go.cipher.host/pkgdex/internal/database"
	"go.cipher.host/pkgdex/internal/meta"
	"go.cipher.host/pkgdex/internal/metric"
	"go.cipher.host/pkgdex/internal/rss"
	"go.cipher.host/pkgdex/internal/search"
	"go.cipher.host/pkgdex/internal/server/handler"
	"go.cipher.host/pkgdex/internal/server/middleware"
	"go.cipher.host/pkgdex/internal/sitemap"
	"go.cipher.host/pkgdex/internal/version"
	"go.cipher.host/pkgdex/internal/wayback"
	"go.cipher.host/pkgdex/static"
)

const (
	// ErrParseHTMLTemplate is the error message returned when the HTML
	// templates can't be parsed.
	ErrParseHTMLTemplate cmdkit.Error = "failed to parse HTML templates; template files may be missing or contain syntax errors"

	// ErrTLSCertificates is the error message returned when the TLS
	// certificates can't be loaded.
	ErrTLSCertificates cmdkit.Error = "failed to load TLS certificate; ensure certificate and key files exist and are valid"

	// ErrCreateSearchManager is the error message returned when the search
	// manager can't be created.
	ErrCreateSearchManager cmdkit.Error = "failed to initialize search functionality"

	// ErrGenerateSitemap is the error message returned when the sitemap can't
	// be generated.
	ErrGenerateSitemap cmdkit.Error = "failed to generate sitemap"

	// ErrGenerateFeed is the error message returned when the RSS feed can't be
	// generated.
	ErrGenerateFeed cmdkit.Error = "failed to generate RSS feed"

	// ErrTrackPackageVersion is the error message returned when the package
	// version can't be tracked.
	ErrTrackPackageVersion cmdkit.Error = "failed to track package version"

	// ErrIndexPackage is the error message returned when the package can't be
	// added to the search index.
	ErrIndexPackage cmdkit.Error = "failed to add package to search index"

	// ErrArchivePackage is the error message returned when the package can't be
	// saved to the Wayback Machine archive.
	ErrArchivePackage cmdkit.Error = "failed to archive package to the Wayback Machine"

	// ErrForcedServerShutdown is the error message returned when the server
	// instance can't be stopped by force.
	ErrForcedServerShutdown cmdkit.Error = "emergency server shutdown failed; server may need manual intervention"

	// ErrCloseDatabase is the error message returned when the database can't be
	// closed.
	ErrCloseDatabase cmdkit.Error = "failed to close server database"
)

// Instance represents a server instance.
type Instance struct {
	cfg        *config.Config
	db         *database.Store
	logger     *slog.Logger
	httpServer *http.Server
	tmpl       *template.Template
}

// NewInstance creates a new server Instance.
func NewInstance(ctx context.Context, cfg *config.Config, db *database.Store, logger *slog.Logger) (*Instance, error) {
	tmpl, err := static.Parse()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrParseHTMLTemplate, err)
	}

	cert, err := tls.X509KeyPair(cfg.Server.TLS.Certificate, cfg.Server.TLS.Key)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrTLSCertificates, err)
	}

	tlsConfig := xtls.ModernServerConfig()
	tlsConfig.Certificates = []tls.Certificate{cert}

	middlewares := []func(http.Handler) http.Handler{
		func(h http.Handler) http.Handler { return xmiddleware.PanicRecovery(logger, h) },
		func(h http.Handler) http.Handler { return xmiddleware.UserAgent(logger, h) },
		func(h http.Handler) http.Handler {
			return xmiddleware.AcceptRequests(
				[]string{
					http.MethodGet,
					http.MethodHead,
					http.MethodOptions,
				},
				logger,
				h,
			)
		},
		middleware.CSP,
	}

	var (
		metrics    = metric.NewStore(db)
		versionMgr = version.NewManager(db)
		archive    = wayback.NewArchive(nil)
	)

	searchMgr, err := search.NewManager(cfg.Database.IndexPath)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCreateSearchManager, err)
	}

	for _, pkg := range cfg.Packages {
		if err = versionMgr.TrackPackage(pkg, time.Now()); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrTrackPackageVersion, err)
		}

		if err = searchMgr.Index(pkg, time.Now()); err != nil {
			return nil, fmt.Errorf("%w (%s): %w", ErrIndexPackage, pkg.Name, err)
		}
	}

	if cfg.Service.Archive {
		if err = archive.SavePackages(ctx, cfg); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrArchivePackage, err)
		}
	}

	sitemapCache, err := sitemap.NewCache(cfg, versionMgr, time.Now())
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGenerateSitemap, err)
	}

	feedCache, err := rss.NewCache(cfg, versionMgr)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGenerateFeed, err)
	}

	var (
		rootHandler      = handler.NewRootHandler(logger, cfg.Service, cfg.Packages, tmpl)
		packageHandler   = handler.NewPackageHandler(logger, metrics, cfg.Packages, cfg.Service, tmpl)
		privacyHandler   = handler.NewPrivacyHandler(logger, cfg.Service, tmpl)
		searchHandler    = handler.NewSearchHandler(logger, searchMgr, cfg.Service, tmpl)
		staticHandler    = handler.NewStaticHandler(logger, tmpl, cfg.Service, static.Files())
		sitemapHandler   = handler.NewSitemapHandler(logger, sitemapCache)
		feedHandler      = handler.NewFeedHandler(logger, feedCache)
		healthHandler    = handler.NewHealthHandler(logger)
		downloadsHandler = handler.NewDownloadsHandler(logger, metrics, cfg.Packages)
	)

	stylePaths, err := static.GetStylePaths()
	if err != nil {
		return nil, fmt.Errorf("getting style paths: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("GET /", xmiddleware.Chain(rootHandler, middlewares...))
	mux.Handle("GET /favicon.ico", xmiddleware.Chain(staticHandler, middlewares...))
	mux.Handle("GET /robots.txt", xmiddleware.Chain(staticHandler, middlewares...))
	mux.Handle("GET /site.webmanifest", xmiddleware.Chain(staticHandler, middlewares...))
	mux.Handle("GET /sitemap.xml", xmiddleware.Chain(sitemapHandler, middlewares...))
	mux.Handle("GET /feed.xml", xmiddleware.Chain(feedHandler, middlewares...))
	mux.Handle("GET /"+stylePaths["main.css"], xmiddleware.Chain(staticHandler, middlewares...))
	mux.Handle("GET /"+stylePaths["highlight.css"], xmiddleware.Chain(staticHandler, middlewares...))
	mux.Handle("GET /assets/img/android-chrome-192x192.png", xmiddleware.Chain(staticHandler, middlewares...))
	mux.Handle("GET /assets/img/android-chrome-512x512.png", xmiddleware.Chain(staticHandler, middlewares...))
	mux.Handle("GET /assets/img/apple-touch-icon.png", xmiddleware.Chain(staticHandler, middlewares...))
	mux.Handle("GET /assets/img/favicon-16x16.png", xmiddleware.Chain(staticHandler, middlewares...))
	mux.Handle("GET /assets/img/favicon-32x32.png", xmiddleware.Chain(staticHandler, middlewares...))
	mux.Handle("GET /assets/img/favicon.svg", xmiddleware.Chain(staticHandler, middlewares...))

	for _, pkg := range cfg.Packages {
		mux.Handle("GET /"+pkg.Name, xmiddleware.Chain(packageHandler, middlewares...))
	}

	mux.Handle("GET /privacy", xmiddleware.Chain(privacyHandler, middlewares...))
	mux.Handle("GET /search", xmiddleware.Chain(searchHandler, middlewares...))

	metaMiddlewares := middlewares
	metaMiddlewares = append(
		metaMiddlewares,
		func(h http.Handler) http.Handler { return middleware.Authorization(cfg.Service.Key, logger, h) },
	)

	mux.Handle("GET /meta/health", xmiddleware.Chain(healthHandler, metaMiddlewares...))
	mux.Handle("GET /meta/downloads", xmiddleware.Chain(downloadsHandler, metaMiddlewares...))
	mux.Handle("GET /meta/debug/pprof/", xmiddleware.Chain(http.HandlerFunc(pprof.Index), metaMiddlewares...))
	mux.Handle("GET /meta/debug/pprof/profile", xmiddleware.Chain(http.HandlerFunc(pprof.Profile), metaMiddlewares...))

	httpServer := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      mux,
		TLSConfig:    tlsConfig,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout),
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout),
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	return &Instance{
		cfg:        cfg,
		db:         db,
		logger:     logger,
		httpServer: httpServer,
		tmpl:       tmpl,
	}, nil
}

// Start starts the server instance.
func (i *Instance) Start() error {
	var key string

	if i.cfg.Service.GeneratedKey {
		key = i.cfg.Service.Key
	} else {
		key = "REDACTED"
	}

	i.logger.Info("initializing server instance",
		slog.String("address", i.cfg.Server.Address),
		slog.String("version", meta.Version),
		slog.String("key", key),
		slog.Int("packages", len(i.cfg.Packages)),
	)

	if err := i.httpServer.ListenAndServeTLS("", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// Stop gracefully stops the server instance.
func (i *Instance) Stop(ctx context.Context) error {
	i.logger.Info("initiating server shutdown")

	if err := i.httpServer.Shutdown(ctx); err != nil {
		i.logger.Error("graceful shutdown failed",
			slog.Any("error", err),
		)

		if err = i.httpServer.Close(); err != nil {
			return fmt.Errorf("%w: %w", ErrForcedServerShutdown, err)
		}
	}

	if err := i.db.Close(); err != nil {
		return fmt.Errorf("%w: %w", ErrCloseDatabase, err)
	}

	i.logger.Info("server shutdown completed successfully")

	return nil
}
