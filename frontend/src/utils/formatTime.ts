export function formatToTimeOfDay(time: number): [number, number, number] {
    const hours = Math.floor(time / 3600);
    const minutes = Math.floor((time % 3600) / 60);
    const seconds = time % 60;

    return [hours, minutes, seconds];
}


export function formatToDate(time: number): string {
    //js takes in time in miliseconds but backend returns in seconds
    const jsTime = new Date(time * 1000);

    return new Intl.DateTimeFormat('en-US', {
        year: 'numeric',
        month: 'short',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        timeZoneName: 'short'
    }).format(jsTime);
}
