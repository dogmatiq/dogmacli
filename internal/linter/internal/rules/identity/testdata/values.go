package subject

import (
	"strings"

	"github.com/dogmatiq/dogma"
)

type (
	namedConstants struct{}
	nonConstant    struct{}
	invalidValues  struct{}
	nonUUIDKey     struct{}
)

// Configure calls c.Identity() named constants in order to verify that the the
// rule does not fail if the values are expressed as literals.
func (namedConstants) Configure(c dogma.ApplicationConfigurer) {
	const name = "name"
	const key = "c82fa717-5217-4a57-84d3-48c58545fb66"

	c.Identity(
		name,
		key,
	)
}

// Configure calls c.Identity() with non-constant values in order to verify that
// the rule does not fail if the values are not known at compile time.
func (nonConstant) Configure(c dogma.ApplicationConfigurer) {
	c.Identity(
		strings.ToUpper("name"),
		strings.ToUpper("key"),
	)
}

// Configure calls c.Identity() with invalid values. It uses validation from
// configkit.NewIdentity().
func (invalidValues) Configure(c dogma.ApplicationConfigurer) {
	c.Identity(
		"", // ruletest: [error] invalid name "", names must be non-empty, printable UTF-8 strings with no whitespace
		"", // ruletest: [error] invalid key "", keys must be RFC 4122 UUIDs
	)
}

// Configure calls c.Identity() with a non-UUID key.
func (nonUUIDKey) Configure(c dogma.ApplicationConfigurer) {
	c.Identity(
		"name",
		"key", // ruletest: [error] invalid key "key", keys must be RFC 4122 UUIDs
	)
}
