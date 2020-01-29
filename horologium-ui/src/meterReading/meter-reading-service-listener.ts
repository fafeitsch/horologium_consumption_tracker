import {MeterReading} from './meter-reading';

export interface MeterReadingServiceListener {
    meterReadingAdded(meterRading: MeterReading): void;
}
