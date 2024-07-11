package guidance1

// Notes:
// Guidance is static data and an access policy, based on rate of change and immediacy, should be determined.

// So far, there are these access policies:
// 1. Global - system calendar, slow changing, the changes are urgent
// 2. Configuration - configurable changing, data can be stale for a period of time.

// Need to create an ops agent for scheduling

// Who controls the freshness/urgency of the data? If the data needs to be very fresh, then a
// constant client polling model is used. Delays on the client will lead to stale data
// If the data doesn't change frequently, and stale data is OK, then the client can poll on intervals.
