# ubk-silo

## Duration of a Subscription

PnYnMnDTnHnMnS
PnW
P<date>T<time>

Durations are a component of time intervals and define the amount of intervening time in a time interval. They should only be used as part of a time interval as prescribed by the standard. Time intervals are discussed in the next section.

Durations are represented by the format P[n]Y[n]M[n]DT[n]H[n]M[n]S or P[n]W as shown to the right. In these representations, the [n] is replaced by the value for each of the date and time elements that follow the [n]. Leading zeros are not required, but the maximum number of digits for each element should be agreed to by the communicating parties. The capital letters P, Y, M, W, D, T, H, M, and S are designators for each of the date and time elements and are not replaced.

    P is the duration designator (historically called "period") placed at the start of the duration representation.
    Y is the year designator that follows the value for the number of years.
    M is the month designator that follows the value for the number of months.
    W is the week designator that follows the value for the number of weeks.
    D is the day designator that follows the value for the number of days.
    T is the time designator that precedes the time components of the representation.
    H is the hour designator that follows the value for the number of hours.
    M is the minute designator that follows the value for the number of minutes.
    S is the second designator that follows the value for the number of seconds.

For example, "P3Y6M4DT12H30M5S" represents a duration of "three years, six months, four days, twelve hours, thirty minutes, and five seconds".

Date and time elements including their designator may be omitted if their value is zero, and lower order elements may also be omitted for reduced precision. For example, "P23DT23H" and "P4Y" are both acceptable duration representations.

To resolve ambiguity, "P1M" is a one-month duration and "PT1M" is a one-minute duration (note the time designator, T, that precedes the time value). The smallest value used may also have a decimal fraction, as in "P0.5Y" to indicate half a year. This decimal fraction may be specified with either a comma or a full stop, as in "P0,5Y" or "P0.5Y". The standard does not prohibit date and time values in a duration representation from exceeding their "carry over points" except as noted below. Thus, "PT36H" could be used as well as "P1DT12H" for representing the same duration. But keep in mind that "PT36H" is not the same as "P1DT12H" when switching from or to Daylight saving time.

Alternatively, a format for duration based on combined date and time representations may be used by agreement between the communicating parties either in the basic format PYYYYMMDDThhmmss or in the extended format P[YYYY]-[MM]-[DD]T[hh]:[mm]:[ss]. For example, the first duration shown above would be "P0003-06-04T12:30:05". However, individual date and time values cannot exceed their moduli (e.g. a value of 13 for the month or 25 for the hour would not be permissible).[18]

* Queries
	* Types

* Updates
* Subscriptionsjo
hola
hola
hola
hola
