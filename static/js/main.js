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
    addUserLocalTime("fp-user-time", fpUserTime);
    addUserLocalTime("sp-user-time", spUserTime);
    addUserLocalTime("tp-user-time", tpUserTime);
    addUserLocalTime("q-user-time", qUserTime);
    if (sprintWeekend) {
        addUserLocalTime("s-user-time", sUserTime);
    }
    addUserLocalTime("r-user-time", rUserTime);

    // Reformat the UTC time to match the users time
    addUtcTime("fp-utc-time", fpUserTime);
    addUtcTime("sp-utc-time", spUserTime);
    addUtcTime("tp-utc-time", tpUserTime);
    addUtcTime("q-utc-time", qUserTime);
    if (sprintWeekend) {
        addUtcTime("s-utc-time", sUserTime);
    }
    addUtcTime("r-utc-time", rUserTime);

    // Register time until updates 
    registerTimeUntil("fp-starts-in", fpUserTime);
    registerTimeUntil("sp-starts-in", spUserTime);
    registerTimeUntil("tp-starts-in", tpUserTime);
    registerTimeUntil("q-starts-in", qUserTime);
    if (sprintWeekend) {
        registerTimeUntil("s-starts-in", sUserTime);
    }
    registerTimeUntil("r-starts-in", rUserTime);

    // Dump the detected timezone on screen so people can self-verify
    document.getElementById("detected-timezone").innerHTML = "Detected timezone: " + timeZoneNameFormat.format(fpUserTime).split(", ")[1];

    // Show all of the elements we hide if JS isn't available
    var jsReqEls = document.getElementsByClassName("js-req");
    while (jsReqEls.length > 0) {
        jsReqEls[0].classList.remove("js-req");
    }
});

function addUserLocalTime(elId, userTime) {
    document.getElementById(elId).innerHTML = weekdayFormat.format(userTime) + " " + localFormat.format(userTime);
}

function addUtcTime(elId, utcTime) {
    document.getElementById(elId).innerHTML = utcFormat.format(utcTime);
}

function registerTimeUntil(elId, countDownDate) {
    // Call the render function right away and then register it to run every second
    renderTimeUntil(elId, countDownDate);
    setInterval(renderTimeUntil, 1000, elId, countDownDate);
}

function renderTimeUntil(elId, countDownDate) {
    // Get the element using elId
    var element = document.getElementById(elId)
    
    // Get today's date and time
    var now = new Date().getTime();

    // Find the distance between now and the count down date
    var distance = countDownDate - now;

    // Set to "Started" and return
    if (distance < 0) {
        element.innerHTML = "Started";
        return;
    }

    // Time calculations for days, hours, minutes and seconds
    var days = Math.floor(distance / (1000 * 60 * 60 * 24));
    var hours = Math.floor((distance % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
    var minutes = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60));
    var seconds = Math.floor((distance % (1000 * 60)) / 1000);

    // Workaround for minus numbers being off by 1
    if (days < 0) {
        days = days + 1;
    }
    if (hours < 0) {
        hours = hours + 1;
    }
    if (minutes < 0) {
        minutes = minutes + 1;
    }
    if (seconds < 0) {
        seconds = seconds + 1;
    }

    // Add a plus to the start if days > 99
    if (days > 99) {
        days = "+99";
    } else {
        days = days.toString().padStart(2, 0);
    }
    
    // Pad the count down to 2 digits each
    hours = hours.toString().padStart(2, 0);
    minutes = minutes.toString().padStart(2, 0);
    seconds = seconds.toString().padStart(2, 0);

    // Update the element with the result
    element.innerHTML = days + "d " + hours + "h " + minutes + "m " + seconds + "s ";
}

function isValidDate(d) {
  return d instanceof Date && !isNaN(d);
}
