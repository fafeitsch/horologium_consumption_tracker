import {MeterReading} from './meter-reading';

export interface MeterReadingServiceListener {
    meterReadingAdded(meterReading: MeterReading): void;
    meterReadingChanged(meterReading: MeterReading): void;
}
