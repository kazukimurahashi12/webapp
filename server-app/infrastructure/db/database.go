package db

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DBマネージャー構造体
type DBManager struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

// シングルトンインスタンス
var dbManager *DBManager

// RDB接続設定情報構造体
// カプセル化
type config struct {
	User       string
	Password   string
	Host       string
	Database   string
	RetryCount int
}

// 環境変数から設定を読み込む
func loadConfig(logger *zap.Logger) (*config, error) {

	// プロジェクトルートディレクトリを取得
	rootDir := os.Getenv("PROJECT_ROOT")
	if rootDir == "" {
		return nil, fmt.Errorf("PROJECT_ROOT environment variable is not set")
	}

	// 環境変数ファイルのパス
	envPath := filepath.Join(rootDir, "build", "db", "data", ".env")

	// 環境変数ファイルを読み込み
	if err := godotenv.Load(envPath); err != nil {
		logger.Warn("Failed to load .env file, using environment variables",
			zap.String("path", envPath),
			zap.Error(err))
	}

	// 必須環境変数のチェック
	requiredEnvVars := []string{"MYSQL_USER", "MYSQL_PASSWORD", "MYSQL_DATABASE"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			return nil, fmt.Errorf("required environment variable is missing: %s", envVar)
		}
	}

	// ホスト設定
	dbHost := os.Getenv("MYSQL_LOCAL_HOST")
	if os.Getenv("DOCKER_ENV") == "true" {
		dbHost = os.Getenv("MYSQL_DOCKER_HOST")
		if dbHost == "" {
			return nil, fmt.Errorf("DOCKER_ENV is true but MYSQL_DOCKER_HOST is not set")
		}
	}

	// リトライ回数設定
	retryCountStr := os.Getenv("RETRYL_COUNT")
	retryCount, err := strconv.Atoi(retryCountStr)
	if err != nil {
		logger.Error("Invalid or missing RETRYL_COUNT value: %v, using default value", zap.Error(err))
		retryCount = 5 // エラー時デフォルト値
	}

	return &config{
		User:       os.Getenv("MYSQL_USER"),
		Password:   os.Getenv("MYSQL_PASSWORD"),
		Host:       dbHost,
		Database:   os.Getenv("MYSQL_DATABASE"),
		RetryCount: retryCount,
	}, nil
}

// DBマネージャーを初期化
func NewDBManager(logger *zap.Logger) *DBManager {
	if dbManager != nil {
		return dbManager
	}

	dbManager = &DBManager{
		DB:     nil,
		Logger: logger,
	}

	return dbManager
}

// データベース接続
func (m *DBManager) Connect() error {
	// 既に接続済みなら何もしない
	if m.DB != nil {
		return nil
	}
	// 設定読み込み
	config, err := loadConfig(m.Logger)
	if err != nil {
		m.Logger.Error("Failed to load configuration", zap.Error(err))
		return err
	}

	// MySQL設定構造体を使用
	cfg := mysqlDriver.Config{
		User:                 config.User,
		Passwd:               config.Password,
		Net:                  "tcp",
		Addr:                 config.Host,
		DBName:               config.Database,
		ParseTime:            true,
		Collation:            "utf8mb4_general_ci",
		AllowNativePasswords: true,
		Loc:                  time.Local,
	}

	// DSN文字列を生成
	dsn := cfg.FormatDSN()
	dialector := mysql.Open(dsn)
	db := connectWithRetry(dialector, config.RetryCount, m.Logger)

	if db == nil {
		return fmt.Errorf("failed to connect to database")
	}

	// DB接続を保持
	m.DB = db

	return nil
}

// 接続リトライ処理
func connectWithRetry(dialector gorm.Dialector, maxRetries int, logger *zap.Logger) *gorm.DB {
	var db *gorm.DB
	var err error

	for retries := 0; retries <= maxRetries; retries++ {
		db, err = gorm.Open(dialector, &gorm.Config{})
		// 接続成功
		if err == nil {
			logger.Info("Connected to database",
				zap.String("dialect", dialector.Name()))
			return db
		}

		if retries < maxRetries {
			retryDelay := time.Second * 2
			logger.Info("Retrying database connection...",
				zap.Int("attempt", retries+1),
				zap.Int("maxRetries", maxRetries),
				zap.Duration("delay", retryDelay))
			time.Sleep(retryDelay)
		} else {
			logger.Error("Failed to connect to database after retries",
				zap.Int("attempts", maxRetries+1),
				zap.Error(err))
		}
	}
	// すべてのリトライ失敗
	logger.Error("Failed to connect to database after retries")
	return nil
}

// DB構造体返却
func (m *DBManager) GetDB() *DBManager {
	if m.DB == nil {
		return nil
	}
	return &DBManager{DB: m.DB}
}

// DB接続状態返却
func (m *DBManager) IsClieintInstance() bool {
	return m.DB != nil
}

// DB接続状態チェック
func (m *DBManager) CheckDBConnection() bool {

	// DB接続を確認
	sqlDB, err := m.DB.DB()
	if err != nil {
		m.Logger.Error("Failed to get database connection", zap.Error(err))
		return false
	}

	if err := sqlDB.Ping(); err != nil {
		m.Logger.Error("Database ping failed", zap.Error(err))
		return false
	}
	return true
}
