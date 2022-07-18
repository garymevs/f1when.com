// Set up Date objects using the Z(UTC) time
var fpUserTime = new Date(fp);
var spUserTime = new Date(sp);
var tpUserTime = new Date(tp);
var qUserTime = new Date(q);
var sUserTime = new Date(s);
var rUserTime = new Date(r);

// Check for sprint weekend
var sprintWeekend = isValidDate(sUserTime);

// Formats
var weekdayFormat = new Intl.DateTimeFormat(undefined, {weekday: "long"});
var localFormat = new Intl.DateTimeFormat(undefined, {dateStyle: "short", timeStyle: "short"});
var utcFormat = new Intl.DateTimeFormat(undefined, {timeZone: "UTC", dateStyle: "short", timeStyle: "short"});
var timeZoneNameFormat = new Intl.DateTimeFormat(undefined, {timeZoneName:"long"});

// Wait for the DOM to load
document.addEventListener("DOMContentLoaded", function() {
    // Add users local time to the table cells
    document.getElementById("fp-user-time").innerHTML = weekdayFormat.format(fpUserTime) + " " + localFormat.format(fpUserTime);
    document.getElementById("sp-user-time").innerHTML = weekdayFormat.format(spUserTime) + " " + localFormat.format(spUserTime);
    document.getElementById("tp-user-time").innerHTML = weekdayFormat.format(tpUserTime) + " " + localFormat.format(tpUserTime);
    document.getElementById("q-user-time").innerHTML = weekdayFormat.format(qUserTime) + " " + localFormat.format(qUserTime);
    if (sprintWeekend) {
        document.getElementById("s-user-time").innerHTML = weekdayFormat.format(sUserTime) + " " + localFormat.format(sUserTime);
    }
    document.getElementById("r-user-time").innerHTML = weekdayFormat.format(rUserTime) + " " + localFormat.format(rUserTime);

    // Reformat the UTC time to match the users time
    document.getElementById("fp-utc-time").innerHTML = utcFormat.format(fpUserTime);
    document.getElementById("sp-utc-time").innerHTML = utcFormat.format(spUserTime);
    document.getElementById("tp-utc-time").innerHTML = utcFormat.format(tpUserTime);
    document.getElementById("q-utc-time").innerHTML = utcFormat.format(qUserTime);
    if (sprintWeekend) {
        document.getElementById("s-utc-time").innerHTML = utcFormat.format(sUserTime);
    }
    document.getElementById("r-utc-time").innerHTML = utcFormat.format(rUserTime);

    // Register time until updates 
    registerTimeUntil(document.getElementById("fp-starts-in"), fpUserTime);
    registerTimeUntil(document.getElementById("sp-starts-in"), spUserTime);
    registerTimeUntil(document.getElementById("tp-starts-in"), tpUserTime);
    registerTimeUntil(document.getElementById("q-starts-in"), qUserTime);
    if (sprintWeekend) {
        registerTimeUntil(document.getElementById("s-starts-in"), sUserTime);
    }
    registerTimeUntil(document.getElementById("r-starts-in"), rUserTime);

    // Dump the detected timezone on screen so people can self-verify
    document.getElementById("detected-timezone").innerHTML = "Detected timezone: " + timeZoneNameFormat.format(fpUserTime).split(", ")[1];

    // Show all of the elements we hide if JS isn't available
    var jsReqEls = document.getElementsByClassName("js-req");
    while (jsReqEls.length > 0) {
        jsReqEls[0].classList.remove("js-req");
    }
});

function registerTimeUntil(element, countDownDate) {
    setInterval(function() {
        // Get today's date and time
        var now = new Date().getTime();

        // Find the distance between now and the count down date
        var distance = countDownDate - now;

        // Time calculations for days, hours, minutes and seconds
        var days = Math.floor(distance / (1000 * 60 * 60 * 24));
        var hours = Math.floor((distance % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
        var minutes = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60));
        var seconds = Math.floor((distance % (1000 * 60)) / 1000);

        if (days > 99) {
            days = "+99";
        } else {
            days = days.toString().padStart(2, 0);
        }
        
        hours = hours.toString().padStart(2, 0);
        minutes = minutes.toString().padStart(2, 0);
        seconds = seconds.toString().padStart(2, 0);

        // Display the result in the element with id="demo"
        element.innerHTML = days + "d " + hours + "h "
        + minutes + "m " + seconds + "s ";
    }, 1000);
}

function isValidDate(d) {
  return d instanceof Date && !isNaN(d);
}
