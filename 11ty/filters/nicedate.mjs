import {DateTime} from 'luxon'
export const niceDate = (date) => {
  return DateTime.fromJSDate(date).toLocaleString(DateTime.DATE_MED)
}
