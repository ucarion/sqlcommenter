// Package sqlcommenter provides a Context-based way to add SQLCommenter
// comments to queries run through the database/sql package.
//
// Driver returns a database/sql/driver.Driver that will read SQLCommenter
// attributes carried on contexts from NewContext.
//
// For convenience methods for using this package with OpenTelemetry, see
// otelsqlcommenter.
package sqlcommenter
