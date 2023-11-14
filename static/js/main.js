// Formats
var weekdayFormat = new Intl.DateTimeFormat(undefined, {weekday: "long"});
var localFormat = new Intl.DateTimeFormat(undefined, {dateStyle: "short", timeStyle: "short"});
var utcFormat = new Intl.DateTimeFormat(undefined, {timeZone: "UTC", dateStyle: "short", timeStyle: "short"});
var timeZoneNameFormat = new Intl.DateTimeFormat(undefined, {timeZoneName:"long"});

// Wait for the DOM to load
document.addEventListener("DOMContentLoaded", function() {
    // Validate and setup user times for each sessions
    sessions.forEach(session => {
        console.log(session.session);
        setupTime(session.session, new Date(session.startTime), session.gmtOffset);
    });

    // Dump the detected timezone on screen so people can self-verify
    document.getElementById("detected-timezone").innerHTML = "Detected timezone: " + timeZoneNameFormat.format(new Date()).split(", ")[1];

    // Show all of the elements we hide if JS isn't available
    var jsReqEls = document.getElementsByClassName("js-req");
    while (jsReqEls.length > 0) {
        jsReqEls[0].classList.remove("js-req");
    }
});

// Do all the user time setup if it's valid 
function setupTime(elId, trackTime, gmtOffset) {
    if (isValidDate(trackTime)) {
        // Apply GMT offset
        trackTime = new Date(trackTime + gmtOffset);

        addUserLocalTime(elId + "-user-time", trackTime);
        addUtcTime(elId + "-utc-time", trackTime);
        registerTimeUntil(elId + "-starts-in", trackTime);
    }
}

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

// Date validation
function isValidDate(d) {
  return d instanceof Date && !isNaN(d);
}
