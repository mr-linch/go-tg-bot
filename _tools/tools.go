//go:build tools

package tools

import (
	_ "github.com/vektra/mockery/v2"
	_ "github.com/volatiletech/sqlboiler/v4"
	_ "github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql"
)
