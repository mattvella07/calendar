import {
    getDate,
    getMonth,
    getYear,
    addHours,
    format
} from "date-fns";

export function formatDateForAPI(date) {
    let d = new Date(date),
        yr = d.getFullYear(),
        mo = d.getMonth() + 1,
        day = d.getDate();

    if (mo < 10) {
        mo = "0" + mo;
    }

    if (day < 10) {
        day = "0" + day;
    }

    return `${yr}-${mo}-${day}`;
}

export function createTimeSlots(date) {
    let arr = [];
    for (let x = -1; x < 23; x++) {
        arr = arr.concat([
            format(
                addHours(
                    new Date(
                        getYear(date),
                        getMonth(date),
                        getDate(date),
                        x
                    ),
                    1
                ),
                "ha"
            ),
            ""
        ]);
    }
    return arr;
}
