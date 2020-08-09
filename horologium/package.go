// Package Horologium provides the ability to compute costs and consumptions of
// of entities like power usage, water usage, …
//
// The package is organized around three basic types:
// meter reading, pricing plan, and series. A Series contains
// a slice of meter readings and a list of pricing plans.
//
// A meter reading represents the actual counter a meter shows at a certain date, e.g.
// on 2019-04-10 the power meter showed 163363 kWh.
//
// A pricing plan states the costs of one unit and in which time frame the pricing plan is valid, e.g.
// between 2019-01 and 2019-06 one kWh costs 0.29 ct and the base price per month is 12 Euro.
//
// A series combines both and offers methods for statistics, e.g. how much was the consumption in 2019 and
// how much did it cost.
package horologium
