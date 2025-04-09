package db

import (
	"errors"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

//////////////////////////////////////////////////////////////////////////////////////////
// カバレッジ率HTML出力
// go test -coverprofile=coverage.out ./ && go tool cover -html=coverage.out
// カバレッジ率txtファイル出力
// go test -coverprofile=coverage.out ./ && go tool cover -func=coverage.out -o coverage.txt
//////////////////////////////////////////////////////////////////////////////////////////

// カスタムタイムアウトエラーを作成
type timeoutError struct{}

func (e *timeoutError) Error() string   { return "timeout" }
func (e *timeoutError) Timeout() bool   { return true }
func (e *timeoutError) Temporary() bool { return true }

// モック用のヘルパー関数
func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *zap.Logger) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	logger := zaptest.NewLogger(t)
	return db, mock, logger
}

// godotenvMock wraps the godotenv package for mocking
type godotenvMock struct {
	originalLoad func(...string) error
	mockLoad     func(...string) error
}

func (m *godotenvMock) Load(filenames ...string) error {
	if m.mockLoad != nil {
		return m.mockLoad(filenames...)
	}
	return m.originalLoad(filenames...)
}

var godotenvMockInstance = &godotenvMock{
	originalLoad: godotenv.Load,
	mockLoad:     func(filenames ...string) error { return nil },
}

// connectWithRetryMock manages mocking of the connectWithRetry function
type connectWithRetryMock struct {
	originalFunc func(gorm.Dialector, int, *zap.Logger) *gorm.DB
	mockFunc     func(gorm.Dialector, int, *zap.Logger) *gorm.DB
}

func (m *connectWithRetryMock) Connect(dialector gorm.Dialector, maxRetries int, logger *zap.Logger) *gorm.DB {
	if m.mockFunc != nil {
		return m.mockFunc(dialector, maxRetries, logger)
	}
	return m.originalFunc(dialector, maxRetries, logger)
}

var connectWithRetryMockInstance = &connectWithRetryMock{
	originalFunc: connectWithRetry,
}

// mockGodotenv enables/disables the godotenv mock
func mockGodotenv(enabled bool) {
	if enabled {
		godotenvMockInstance.mockLoad = func(filenames ...string) error { return nil }
	} else {
		godotenvMockInstance.mockLoad = nil
	}
}

// mockConnectWithRetry sets up the mock for connectWithRetry
func mockConnectWithRetry(db *gorm.DB, success bool) {
	if success {
		connectWithRetryMockInstance.mockFunc = func(dialector gorm.Dialector, maxRetries int, logger *zap.Logger) *gorm.DB {
			return db
		}
	} else {
		connectWithRetryMockInstance.mockFunc = func(dialector gorm.Dialector, maxRetries int, logger *zap.Logger) *gorm.DB {
			return nil
		}
	}
}

// resetConnectWithRetry resets the mock
func resetConnectWithRetry() {
	connectWithRetryMockInstance.mockFunc = nil
}

// IsClientInstance checks if DB connection exists
func (m *DBManager) IsClientInstance() bool {
	return m.db != nil
}

// NewDBManagerのテスト
func TestNewDBManager(t *testing.T) {
	t.Run("SuccessStandardFlow", func(t *testing.T) {
		// Arrange ---
		logger := zaptest.NewLogger(t)
		// シングルトンインスタンスをリセット
		dbManager = nil

		// Act ---
		manager := NewDBManager(logger)

		// Assert ---
		assert.NotNil(t, manager)
		assert.Nil(t, manager.db)
		assert.Equal(t, logger, manager.logger)
	})

	t.Run("AlreadyInitialized", func(t *testing.T) {
		// Arrange ---
		logger := zaptest.NewLogger(t)
		// 既に初期化されているシングルトンインスタンスを使用
		firstManager := NewDBManager(logger)

		// Act ---
		secondManager := NewDBManager(logger)

		// Assert ---
		assert.Equal(t, firstManager, secondManager)
	})
}

// loadConfigのテスト
func TestLoadConfig(t *testing.T) {
	t.Run("正常系（ローカル環境）", func(t *testing.T) {
		// Arrange ---
		logger := zaptest.NewLogger(t)
		mockGodotenv(true)
		defer mockGodotenv(false)

		// 環境変数を設定
		envVars := map[string]string{
			"PROJECT_ROOT":     "/test/root",
			"MYSQL_USER":       "testuser",
			"MYSQL_PASSWORD":   "testpass",
			"MYSQL_DATABASE":   "testdb",
			"MYSQL_LOCAL_HOST": "localhost:3306",
			"RETRYL_COUNT":     "3",
		}

		originalEnv := make(map[string]string)
		for k := range envVars {
			originalEnv[k] = os.Getenv(k)
		}

		for k, v := range envVars {
			os.Setenv(k, v)
		}
		defer func() {
			// 環境変数を元に戻す
			for k, v := range originalEnv {
				os.Setenv(k, v)
			}
		}()

		expectedConfig := &config{
			User:       "testuser",
			Password:   "testpass",
			Host:       "localhost:3306",
			Database:   "testdb",
			RetryCount: 3,
		}

		// Act ---
		cfg, err := loadConfig(logger)

		// Assert ---
		assert.NoError(t, err)
		assert.Equal(t, expectedConfig, cfg)
	})

	t.Run("正常系（Docker環境）", func(t *testing.T) {
		// Arrange ---
		logger := zaptest.NewLogger(t)
		mockGodotenv(true)
		defer mockGodotenv(false)

		// 環境変数を設定
		envVars := map[string]string{
			"PROJECT_ROOT":      "/test/root",
			"MYSQL_USER":        "testuser",
			"MYSQL_PASSWORD":    "testpass",
			"MYSQL_DATABASE":    "testdb",
			"MYSQL_LOCAL_HOST":  "localhost:3306",
			"MYSQL_DOCKER_HOST": "mysql:3306",
			"DOCKER_ENV":        "true",
			"RETRYL_COUNT":      "3",
		}

		originalEnv := make(map[string]string)
		for k := range envVars {
			originalEnv[k] = os.Getenv(k)
		}

		for k, v := range envVars {
			os.Setenv(k, v)
		}
		defer func() {
			// 環境変数を元に戻す
			for k, v := range originalEnv {
				os.Setenv(k, v)
			}
		}()

		expectedConfig := &config{
			User:       "testuser",
			Password:   "testpass",
			Host:       "mysql:3306",
			Database:   "testdb",
			RetryCount: 3,
		}

		// Act ---
		cfg, err := loadConfig(logger)

		// Assert ---
		assert.NoError(t, err)
		assert.Equal(t, expectedConfig, cfg)
	})

	t.Run("PROJECT_ROOT未設定エラー", func(t *testing.T) {
		// Arrange ---
		logger := zaptest.NewLogger(t)
		mockGodotenv(true)
		defer mockGodotenv(false)

		// 環境変数を設定
		envVars := map[string]string{
			"PROJECT_ROOT": "",
		}

		originalEnv := make(map[string]string)
		for k := range envVars {
			originalEnv[k] = os.Getenv(k)
		}

		for k, v := range envVars {
			os.Setenv(k, v)
		}
		defer func() {
			// 環境変数を元に戻す
			for k, v := range originalEnv {
				os.Setenv(k, v)
			}
		}()

		// Act ---
		cfg, err := loadConfig(logger)

		// Assert ---
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "PROJECT_ROOT environment variable is not set")
		assert.Nil(t, cfg)
	})

	t.Run("必須環境変数不足エラー", func(t *testing.T) {
		// Arrange ---
		logger := zaptest.NewLogger(t)
		mockGodotenv(true)
		defer mockGodotenv(false)

		// 環境変数を設定
		envVars := map[string]string{
			"PROJECT_ROOT":   "/test/root",
			"MYSQL_USER":     "testuser",
			"MYSQL_PASSWORD": "",
			"MYSQL_DATABASE": "testdb",
		}

		originalEnv := make(map[string]string)
		for k := range envVars {
			originalEnv[k] = os.Getenv(k)
		}

		for k, v := range envVars {
			os.Setenv(k, v)
		}
		defer func() {
			// 環境変数を元に戻す
			for k, v := range originalEnv {
				os.Setenv(k, v)
			}
		}()

		// Act ---
		cfg, err := loadConfig(logger)

		// Assert ---
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "required environment variable is missing: MYSQL_PASSWORD")
		assert.Nil(t, cfg)
	})

	t.Run("Docker環境でホスト未設定エラー", func(t *testing.T) {
		// Arrange ---
		logger := zaptest.NewLogger(t)
		mockGodotenv(true)
		defer mockGodotenv(false)

		// 環境変数を設定
		envVars := map[string]string{
			"PROJECT_ROOT":      "/test/root",
			"MYSQL_USER":        "testuser",
			"MYSQL_PASSWORD":    "testpass",
			"MYSQL_DATABASE":    "testdb",
			"MYSQL_LOCAL_HOST":  "localhost:3306",
			"MYSQL_DOCKER_HOST": "",
			"DOCKER_ENV":        "true",
		}

		originalEnv := make(map[string]string)
		for k := range envVars {
			originalEnv[k] = os.Getenv(k)
		}

		for k, v := range envVars {
			os.Setenv(k, v)
		}
		defer func() {
			// 環境変数を元に戻す
			for k, v := range originalEnv {
				os.Setenv(k, v)
			}
		}()

		// Act ---
		cfg, err := loadConfig(logger)

		// Assert ---
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "DOCKER_ENV is true but MYSQL_DOCKER_HOST is not set")
		assert.Nil(t, cfg)
	})

	t.Run("リトライ回数が無効な場合のデフォルト値テスト", func(t *testing.T) {
		// Arrange ---
		logger := zaptest.NewLogger(t)
		mockGodotenv(true)
		defer mockGodotenv(false)

		// 環境変数を設定
		envVars := map[string]string{
			"PROJECT_ROOT":     "/test/root",
			"MYSQL_USER":       "testuser",
			"MYSQL_PASSWORD":   "testpass",
			"MYSQL_DATABASE":   "testdb",
			"MYSQL_LOCAL_HOST": "localhost:3306",
			"RETRYL_COUNT":     "invalid",
		}

		originalEnv := make(map[string]string)
		for k := range envVars {
			originalEnv[k] = os.Getenv(k)
		}

		for k, v := range envVars {
			os.Setenv(k, v)
		}
		defer func() {
			// 環境変数を元に戻す
			for k, v := range originalEnv {
				os.Setenv(k, v)
			}
		}()

		expectedConfig := &config{
			User:       "testuser",
			Password:   "testpass",
			Host:       "localhost:3306",
			Database:   "testdb",
			RetryCount: 5, // デフォルト値
		}

		// Act ---
		cfg, err := loadConfig(logger)

		// Assert ---
		assert.NoError(t, err)
		assert.Equal(t, expectedConfig, cfg)
	})
}

// Connectのテスト
func TestConnect(t *testing.T) {
	t.Run("正常系の接続", func(t *testing.T) {
		// Arrange ---
		logger := zaptest.NewLogger(t)
		db, _, _ := setupMockDB(t)
		mockGodotenv(true)
		defer mockGodotenv(false)

		// 環境変数を設定
		envVars := map[string]string{
			"PROJECT_ROOT":     "/test/root",
			"MYSQL_USER":       "testuser",
			"MYSQL_PASSWORD":   "testpass",
			"MYSQL_DATABASE":   "testdb",
			"MYSQL_LOCAL_HOST": "localhost:3306",
			"RETRYL_COUNT":     "3",
		}

		originalEnv := make(map[string]string)
		for k := range envVars {
			originalEnv[k] = os.Getenv(k)
		}

		for k, v := range envVars {
			os.Setenv(k, v)
		}
		defer func() {
			// 環境変数を元に戻す
			for k, v := range originalEnv {
				os.Setenv(k, v)
			}
		}()

		// connectWithRetryをモック化
		mockConnectWithRetry(db, true)
		defer resetConnectWithRetry()

		// DBManagerインスタンスを作成
		manager := &DBManager{
			logger: logger,
		}

		// Act ---
		err := manager.Connect()

		// Assert ---
		assert.NoError(t, err)
		assert.NotNil(t, manager.db)
	})

	t.Run("既に接続済みの場合", func(t *testing.T) {
		// Arrange ---
		logger := zaptest.NewLogger(t)
		db, _, _ := setupMockDB(t)
		mockGodotenv(true)
		defer mockGodotenv(false)

		// 環境変数を設定
		envVars := map[string]string{
			"PROJECT_ROOT":     "/test/root",
			"MYSQL_USER":       "testuser",
			"MYSQL_PASSWORD":   "testpass",
			"MYSQL_DATABASE":   "testdb",
			"MYSQL_LOCAL_HOST": "localhost:3306",
			"RETRYL_COUNT":     "3",
		}

		originalEnv := make(map[string]string)
		for k := range envVars {
			originalEnv[k] = os.Getenv(k)
		}

		for k, v := range envVars {
			os.Setenv(k, v)
		}
		defer func() {
			// 環境変数を元に戻す
			for k, v := range originalEnv {
				os.Setenv(k, v)
			}
		}()

		// connectWithRetryをモック化
		mockConnectWithRetry(db, true)
		defer resetConnectWithRetry()

		// DBManagerインスタンスを作成
		manager := &DBManager{
			logger: logger,
			db:     db, // 既に接続済み
		}

		// Act ---
		err := manager.Connect()

		// Assert ---
		assert.NoError(t, err)
		assert.NotNil(t, manager.db)
	})

	t.Run("設定読み込みエラー", func(t *testing.T) {
		// Arrange ---
		logger := zaptest.NewLogger(t)
		db, _, _ := setupMockDB(t)
		mockGodotenv(true)
		defer mockGodotenv(false)

		// 環境変数を設定
		envVars := map[string]string{
			"PROJECT_ROOT": "", // エラーになる
		}

		originalEnv := make(map[string]string)
		for k := range envVars {
			originalEnv[k] = os.Getenv(k)
		}

		for k, v := range envVars {
			os.Setenv(k, v)
		}
		defer func() {
			// 環境変数を元に戻す
			for k, v := range originalEnv {
				os.Setenv(k, v)
			}
		}()

		// connectWithRetryをモック化
		mockConnectWithRetry(db, true)
		defer resetConnectWithRetry()

		// DBManagerインスタンスを作成
		manager := &DBManager{
			logger: logger,
		}

		// Act ---
		err := manager.Connect()

		// Assert ---
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "PROJECT_ROOT environment variable is not set")
	})

	t.Run("接続エラー", func(t *testing.T) {
		// Arrange ---
		logger := zaptest.NewLogger(t)
		db, _, _ := setupMockDB(t)
		mockGodotenv(true)
		defer mockGodotenv(false)

		// 環境変数を設定
		envVars := map[string]string{
			"PROJECT_ROOT":     "/test/root",
			"MYSQL_USER":       "testuser",
			"MYSQL_PASSWORD":   "testpass",
			"MYSQL_DATABASE":   "testdb",
			"MYSQL_LOCAL_HOST": "localhost:3306",
			"RETRYL_COUNT":     "3",
		}

		originalEnv := make(map[string]string)
		for k := range envVars {
			originalEnv[k] = os.Getenv(k)
		}

		for k, v := range envVars {
			os.Setenv(k, v)
		}
		defer func() {
			// 環境変数を元に戻す
			for k, v := range originalEnv {
				os.Setenv(k, v)
			}
		}()

		// connectWithRetryをモック化
		mockConnectWithRetry(db, false) // 接続失敗
		defer resetConnectWithRetry()

		// DBManagerインスタンスを作成
		manager := &DBManager{
			logger: logger,
		}

		// Act ---
		err := manager.Connect()

		// Assert ---
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to connect to database")
	})
}

// connectWithRetryのテスト
func TestConnectWithRetry(t *testing.T) {
	t.Run("最初の試行で成功", func(t *testing.T) {
		// Arrange ---
		logger := zaptest.NewLogger(t)
		sqlDB, mock, _ := sqlmock.New()

		// connectWithRetry関数を元に戻す
		resetConnectWithRetry()

		dialector := mysql.New(mysql.Config{
			Conn:                      sqlDB,
			SkipInitializeWithVersion: true,
		})

		// Act ---
		db := connectWithRetry(dialector, 3, logger)

		// Assert ---
		assert.NotNil(t, db)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("すべての試行で失敗", func(t *testing.T) {
		// Arrange ---
		logger := zaptest.NewLogger(t)
		sqlDB, mock, _ := sqlmock.New()

		// connectWithRetry関数を元に戻す
		resetConnectWithRetry()

		// 接続を閉じてエラーを起こす
		sqlDB.Close()

		dialector := mysql.New(mysql.Config{
			Conn:                      sqlDB,
			SkipInitializeWithVersion: true,
		})

		// Act ---
		db := connectWithRetry(dialector, 2, logger)

		// Assert ---
		assert.Nil(t, db)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// GetDBのテスト
func TestGetDB(t *testing.T) {
	t.Run("DB接続がある場合", func(t *testing.T) {
		// Arrange ---
		logger := zaptest.NewLogger(t)
		db, _, _ := setupMockDB(t)

		manager := &DBManager{
			logger: logger,
			db:     db,
		}

		// Act ---
		result := manager.GetDB()

		// Assert ---
		assert.NotNil(t, result)
		assert.NotNil(t, result.db)
	})

	t.Run("DB接続がない場合", func(t *testing.T) {
		// Arrange ---
		logger := zaptest.NewLogger(t)

		manager := &DBManager{
			logger: logger,
			db:     nil,
		}

		// Act ---
		result := manager.GetDB()

		// Assert ---
		assert.Nil(t, result)
	})
}

// IsClientInstanceのテスト
func TestIsClientInstance(t *testing.T) {
	t.Run("DB接続がある場合", func(t *testing.T) {
		// Arrange ---
		logger := zaptest.NewLogger(t)
		db, _, _ := setupMockDB(t)

		manager := &DBManager{
			logger: logger,
			db:     db,
		}

		// Act ---
		result := manager.IsClientInstance()

		// Assert ---
		assert.True(t, result)
	})

	t.Run("DB接続がない場合", func(t *testing.T) {
		// Arrange ---
		logger := zaptest.NewLogger(t)

		manager := &DBManager{
			logger: logger,
			db:     nil,
		}

		// Act ---
		result := manager.IsClientInstance()

		// Assert ---
		assert.False(t, result)
	})
}

// CheckDBConnectionのテスト
func TestCheckDBConnection(t *testing.T) {
	t.Run("接続成功", func(t *testing.T) {
		// Arrange ---
		logger := zaptest.NewLogger(t)
		sqlDB, mock, _ := sqlmock.New()

		// Pingが成功するよう設定
		mock.ExpectPing()

		dialector := mysql.New(mysql.Config{
			Conn:                      sqlDB,
			SkipInitializeWithVersion: true,
		})

		gormDB, err := gorm.Open(dialector, &gorm.Config{})
		require.NoError(t, err)

		manager := &DBManager{
			db:     gormDB,
			logger: logger,
		}

		// Act ---
		result := manager.CheckDBConnection()

		// Assert ---
		assert.True(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Pingエラー", func(t *testing.T) {
		// Arrange ---
		logger := zaptest.NewLogger(t)
		sqlDB, mock, _ := sqlmock.New()

		// Pingがエラーを返すよう設定
		mock.ExpectPing().WillReturnError(errors.New("ping error"))

		dialector := mysql.New(mysql.Config{
			Conn:                      sqlDB,
			SkipInitializeWithVersion: true,
		})

		gormDB, err := gorm.Open(dialector, &gorm.Config{})
		require.NoError(t, err)

		manager := &DBManager{
			db:     gormDB,
			logger: logger,
		}

		// Act ---
		result := manager.CheckDBConnection()

		// Assert ---
		assert.False(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
