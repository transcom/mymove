import { includedStatusesForCalculatingWeights } from 'hooks/custom';
import { shipmentStatuses } from 'constants/shipments';

describe('includedStatusesForCalculatingWeights returns true for approved, diversion requested, or cancellation requested', () => {
  it.each([
    [shipmentStatuses.DRAFT, false],
    [shipmentStatuses.SUBMITTED, false],
    [shipmentStatuses.APPROVED, true],
    [shipmentStatuses.REJECTED, false],
    [shipmentStatuses.CANCELLATION_REQUESTED, true],
    [shipmentStatuses.CANCELED, false],
    [shipmentStatuses.DIVERSION_REQUESTED, true],
    ['FAKE_STATUS', false],
  ])('checks if a shipment with status %s should be included: %b', (status, isIncluded) => {
    expect(includedStatusesForCalculatingWeights(status)).toBe(isIncluded);
  });
});
