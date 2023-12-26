import { renderHook } from '@testing-library/react-hooks';

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
  it('correctly calculates the total weight with diversions present', () => {
    // Lowest diversion is the reweigh found at 600, plus other eligible weights 2000 + 1100 + 800 = a total of 4500
    const mtoShipments = [
      {
        primeActualWeight: 1000,
        reweigh: {
          weight: null,
        },
        status: shipmentStatuses.APPROVED,
        diversion: true,
      },
      {
        primeActualWeight: 1500,
        reweigh: {
          weight: 1300,
        },
        status: shipmentStatuses.APPROVED,
        diversion: true,
      },
      {
        primeActualWeight: 1500,
        reweigh: {
          weight: 600,
        },
        status: shipmentStatuses.APPROVED,
        diversion: true,
      },
      {
        primeActualWeight: 2000,
        reweigh: {
          weight: null,
        },
        status: shipmentStatuses.APPROVED,
      },
      {
        primeActualWeight: 1200,
        reweigh: {
          weight: 1100,
        },
        status: shipmentStatuses.CANCELLATION_REQUESTED,
      },
      {
        primeActualWeight: 800,
        reweigh: {
          weight: null,
        },
        status: shipmentStatuses.DIVERSION_REQUESTED,
      },
    ];

    const { result } = renderHook(() => useCalculatedWeightRequested(mtoShipments));

    expect(result.current).toBe(4500);
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
  it('useCalculatedTotalEstimatedWeight with diversions present', () => {
    // Two approved shipments are now present with an estimated weight of 3500 and 4000
    // Since they are diversions, it should only take the lowest weight
    let mtoShipments = [
      {
        primeEstimatedWeight: 1000,
        calculatedBillableWeight: 10,
        status: shipmentStatuses.DRAFT,
      },
      {
        primeEstimatedWeight: 3500,
        calculatedBillableWeight: 500,
        status: shipmentStatuses.APPROVED,
        diversion: true,
      },
      {
        primeEstimatedWeight: 4000,
        calculatedBillableWeight: 500,
        status: shipmentStatuses.APPROVED,
        diversion: true,
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

    expect(result.current).toBe(5500);

    mtoShipments = mtoShipments.concat([
      { primeEstimatedWeight: 2000, calculatedBillableWeight: 100, status: shipmentStatuses.APPROVED },
    ]);
    rerender();

    expect(result.current).toBe(7500);
  });
});
