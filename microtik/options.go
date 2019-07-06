package microtik

import "strings"

// RouteID is a functional option. It is used when working with ip/routes
// on microtik router. In order to work with routes there needs to be a
// way to identify them. This most reliable way is to identify them through
// their unique comment. With this option a route can be registered with an
// arbitrary name and a string matching the route's comment in the microtik
// router.
func RouteID(name, comment string) Option {
	return func(m *Microtik) {
		name = strings.ToLower(name)
		m.routeIDs[name] = comment
	}
}
