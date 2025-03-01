package config

import (
	"regexp"
	"strings"

	"github.com/photoprism/photoprism/pkg/rnd"
	"golang.org/x/crypto/bcrypt"
)

const (
	AuthModePublic   = "public"
	AuthModePassword = "password"
)

func isBcrypt(s string) bool {
	b, err := regexp.MatchString(`^\$2[ayb]\$.{56}$`, s)
	if err != nil {
		return false
	}
	return b
}

// Public checks if app runs in public mode and requires no authentication.
func (c *Config) Public() bool {
	if c.Demo() {
		return true
	}

	return c.options.Public
}

// SetPublic changes authentication while instance is running, for testing purposes only.
func (c *Config) SetPublic(enabled bool) {
	if c.Debug() {
		c.options.Public = enabled
	}
}

// AdminPassword returns the initial admin password.
func (c *Config) AdminPassword() string {
	return c.options.AdminPassword
}

// AuthMode returns the authentication mode.
func (c *Config) AuthMode() string {
	if c.Public() {
		return AuthModePublic
	} else if m := strings.ToLower(strings.TrimSpace(c.options.AuthMode)); m != "" {
		return m
	}

	return AuthModePassword
}

// Auth checks if authentication is required.
func (c *Config) Auth() bool {
	return c.AuthMode() != AuthModePublic
}

// CheckPassword compares given password p with the admin password
func (c *Config) CheckPassword(p string) bool {
	ap := c.AdminPassword()

	if isBcrypt(ap) {
		err := bcrypt.CompareHashAndPassword([]byte(ap), []byte(p))
		return err == nil
	}

	return ap == p
}

// InvalidDownloadToken checks if the token is invalid.
func (c *Config) InvalidDownloadToken(t string) bool {
	return c.DownloadToken() != t
}

// DownloadToken returns the DOWNLOAD api token (you can optionally use a static value for permanent caching).
func (c *Config) DownloadToken() string {
	if c.options.DownloadToken == "" {
		c.options.DownloadToken = rnd.GenerateToken(8)
	}

	return c.options.DownloadToken
}

// InvalidPreviewToken checks if the preview token is invalid.
func (c *Config) InvalidPreviewToken(t string) bool {
	return c.PreviewToken() != t && c.DownloadToken() != t
}

// PreviewToken returns the preview image api token (based on the unique storage serial by default).
func (c *Config) PreviewToken() string {
	if c.options.PreviewToken == "" {
		if c.Public() {
			c.options.PreviewToken = "public"
		} else if c.Serial() == "" {
			return "********"
		} else {
			c.options.PreviewToken = c.SerialChecksum()
		}
	}

	return c.options.PreviewToken
}
