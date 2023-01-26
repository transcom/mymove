import { renderHook } from '@testing-library/react';

import {
  includedStatusesForCalculatingWeights,
  useCalculatedTotalBillableWeight,
  useCalculatedWeightRequested,
  useCalculatedEstimatedWeight,
} from 'hooks/custom';
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

describe('for all shipments that are approved, have a cancellation requested, or have a diversion requested', () => {
  it('useCalculatedTotalBillableWeight returns the calculated billable weight', () => {
    let mtoShipments = [
      {
        calculatedBillableWeight: 10,
        status: shipmentStatuses.DRAFT,
      },
      {
        calculatedBillableWeight: 500,
        status: shipmentStatuses.APPROVED,
      },
      {
        calculatedBillableWeight: 200,
        status: shipmentStatuses.CANCELLATION_REQUESTED,
      },
      {
        calculatedBillableWeight: 300,
        status: shipmentStatuses.DIVERSION_REQUESTED,
      },
    ];

    const { result, rerender } = renderHook(() => useCalculatedTotalBillableWeight(mtoShipments));

    expect(result.current).toBe(1000);

    mtoShipments = mtoShipments.concat([{ calculatedBillableWeight: 100, status: shipmentStatuses.APPROVED }]);
    rerender();

    expect(result.current).toBe(1100);
  });

  it('useCalculatedWeightRequested returns the calculated billable weight using the lower value between the prime actual weight and the reweigh weight', () => {
    let mtoShipments = [
      {
        primeActualWeight: 10,
        status: shipmentStatuses.DRAFT,
        reweigh: {
          weight: 5,
        },
      },
      {
        primeActualWeight: 2000,
        status: shipmentStatuses.APPROVED,
        reweigh: {
          weight: 300,
        },
      },
      {
        primeActualWeight: 100,
        status: shipmentStatuses.APPROVED,
      },
      {
        primeActualWeight: 1000,
        status: shipmentStatuses.CANCELLATION_REQUESTED,
        reweigh: {
          weight: 200,
        },
      },
      {
        primeActualWeight: 400,
        status: shipmentStatuses.DIVERSION_REQUESTED,
        reweigh: {
          weight: 3000,
        },
      },
    ];

    const { result, rerender } = renderHook(() => useCalculatedWeightRequested(mtoShipments));

    expect(result.current).toBe(1000);

    mtoShipments = mtoShipments.concat([
      { primeActualWeight: 100, status: shipmentStatuses.APPROVED, reweigh: { weight: 3000 } },
    ]);
    rerender();

    expect(result.current).toBe(1100);
  });
  it('useCalculatedTotalEstimatedWeight', () => {
    let mtoShipments = [
      {
        primeEstimatedWeight: 1000,
        calculatedBillableWeight: 10,
        status: shipmentStatuses.DRAFT,
      },
      {
        primeEstimatedWeight: 4000,
        calculatedBillableWeight: 500,
        status: shipmentStatuses.APPROVED,
      },
      {
        primeEstimatedWeight: 1000,
        calculatedBillableWeight: 200,
        status: shipmentStatuses.CANCELLATION_REQUESTED,
      },
      {
        primeEstimatedWeight: 1000,
        calculatedBillableWeight: 300,
        status: shipmentStatuses.DIVERSION_REQUESTED,
      },
    ];

    const { result, rerender } = renderHook(() => useCalculatedEstimatedWeight(mtoShipments));

    expect(result.current).toBe(6000);

    mtoShipments = mtoShipments.concat([
      { primeEstimatedWeight: 2000, calculatedBillableWeight: 100, status: shipmentStatuses.APPROVED },
    ]);
    rerender();

    expect(result.current).toBe(8000);
  });
});
