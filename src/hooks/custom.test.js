import { renderHook } from '@testing-library/react-hooks';

import {
  includedStatusesForCalculatingWeights,
  useCalculatedTotalBillableWeight,
  useCalculatedWeightRequested,
  useCalculatedEstimatedWeight,
} from 'hooks/custom';
import { shipmentStatuses } from 'constants/shipments';

describe('includedStatusesForCalculatingWeights returns true for approved, approvals requested, diversion requested, or cancellation requested', () => {
  it.each([
    [shipmentStatuses.DRAFT, false],
    [shipmentStatuses.SUBMITTED, false],
    [shipmentStatuses.APPROVED, true],
    [shipmentStatuses.APPROVALS_REQUESTED, true],
    [shipmentStatuses.REJECTED, false],
    [shipmentStatuses.CANCELLATION_REQUESTED, true],
    [shipmentStatuses.CANCELED, false],
    [shipmentStatuses.DIVERSION_REQUESTED, true],
    ['FAKE_STATUS', false],
  ])('checks if a shipment with status %s should be included: %b', (status, isIncluded) => {
    expect(includedStatusesForCalculatingWeights(status)).toBe(isIncluded);
  });
});

describe('for all shipments that are approved, approvals requested, have a cancellation requested, or have a diversion requested', () => {
  it('useCalculatedTotalBillableWeight returns the calculated billable weight', () => {
    let mtoShipments = [
      {
        calculatedBillableWeight: 10,
        primeEstimatedWeight: 10,
        primeActualWeight: 10,
        shipmentType: 'HHG',
        status: shipmentStatuses.DRAFT,
      },
      {
        calculatedBillableWeight: 500,
        primeEstimatedWeight: 10,
        primeActualWeight: 10,
        shipmentType: 'HHG',
        status: shipmentStatuses.APPROVED,
      },
      {
        calculatedBillableWeight: 200,
        primeEstimatedWeight: 10,
        primeActualWeight: 10,
        shipmentType: 'HHG',
        status: shipmentStatuses.CANCELLATION_REQUESTED,
      },
      {
        calculatedBillableWeight: 300,
        primeEstimatedWeight: 10,
        primeActualWeight: 10,
        shipmentType: 'HHG',
        status: shipmentStatuses.DIVERSION_REQUESTED,
      },
    ];

    const { result, rerender } = renderHook(() => useCalculatedTotalBillableWeight(mtoShipments));

    expect(result.current).toBe(30);

    mtoShipments = mtoShipments.concat([
      {
        calculatedBillableWeight: 300,
        primeEstimatedWeight: 100,
        primeActualWeight: 100,
        shipmentType: 'HHG',
        status: shipmentStatuses.APPROVED,
      },
    ]);
    rerender();

    expect(result.current).toBe(130);
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
  it('useCalculatedWeightRequested returns the calculated billable weight using the lower value between the prime actual weight and the reweigh weight with a single divesion chain present', () => {
    const mtoShipments = [
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
      {
        primeActualWeight: 400,
        status: shipmentStatuses.APPROVED,
        diversion: true,
        pickupAddress: { city: 'City1' },
        destinationAddress: { city: 'City2' },
        reweigh: {
          weight: 3000,
        },
      },
      {
        primeActualWeight: 400,
        status: shipmentStatuses.APPROVED,
        diversion: true,
        pickupAddress: { city: 'City2' },
        destinationAddress: { city: 'City3' },
        reweigh: {
          weight: 3000,
        },
      },
    ];

    const { result } = renderHook(() => useCalculatedWeightRequested(mtoShipments));

    expect(result.current).toBe(1000 + 400); // Add the lowest actual weight from the diversion chain
  });
  it('useCalculatedWeightRequested returns the calculated billable weight using the lower value between the prime actual weight and the reweigh weight with two divesion chains present', () => {
    const mtoShipments = [
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
      {
        id: 'parent1',
        primeActualWeight: 400,
        status: shipmentStatuses.APPROVED,
        diversion: true,
        pickupAddress: { city: 'City1' },
        destinationAddress: { city: 'City2' },
        reweigh: {
          weight: 3000,
        },
      },
      {
        id: 'child1',
        primeActualWeight: 400,
        status: shipmentStatuses.APPROVED,
        diversion: true,
        pickupAddress: { city: 'City2' },
        destinationAddress: { city: 'City3' },
        reweigh: {
          weight: 3000,
        },
      },
      {
        id: 'parent2',
        primeActualWeight: 800,
        status: shipmentStatuses.APPROVED,
        diversion: true,
        pickupAddress: { city: 'CityA' },
        destinationAddress: { city: 'CityB' },
        reweigh: {
          weight: 3000,
        },
      },
      {
        id: 'child2',
        primeActualWeight: 1000,
        status: shipmentStatuses.APPROVED,
        diversion: true,
        pickupAddress: { city: 'CityB' },
        destinationAddress: { city: 'CityC' },
        reweigh: {
          weight: 3000,
        },
      },
    ];

    const { result } = renderHook(() => useCalculatedWeightRequested(mtoShipments));

    expect(result.current).toBe(1000 + 400 + 800); // Add the lowest actual weight from the two diversion chains
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
    const mtoShipments = [
      {
        primeEstimatedWeight: 1000,
        calculatedBillableWeight: 10,
        status: shipmentStatuses.DRAFT,
      },
      {
        primeEstimatedWeight: 1000,
        calculatedBillableWeight: 200,
        status: shipmentStatuses.CANCELLATION_REQUESTED,
      },
      {
        id: 'parent',
        primeEstimatedWeight: 1000,
        calculatedBillableWeight: 300,
        pickupAddress: { city: 'CityA' },
        destinationAddress: { city: 'CityB' },
        diversion: true,
        status: shipmentStatuses.DIVERSION_REQUESTED,
      },
      {
        id: 'child',
        primeEstimatedWeight: 1500,
        calculatedBillableWeight: 300,
        pickupAddress: { city: 'CityB' },
        destinationAddress: { city: 'CityC' },
        diversion: true,
        status: shipmentStatuses.APPROVED,
      },
    ];

    const { result } = renderHook(() => useCalculatedEstimatedWeight(mtoShipments));

    expect(result.current).toBe(1000 + 1000);
  });
  it('calculates the lowest weight of shipments with 2 parent diversions', () => {
    const mtoShipments = [
      {
        id: 'parent1',
        primeEstimatedWeight: 3000,
        status: shipmentStatuses.APPROVED,
        diversion: true,
        pickupAddress: { city: 'CityA' },
        destinationAddress: { city: 'CityB' },
      },
      {
        id: 'child1',
        primeEstimatedWeight: 2500,
        status: shipmentStatuses.APPROVED,
        diversion: true,
        pickupAddress: { city: 'CityB' },
        destinationAddress: { city: 'CityC' },
      },
      {
        id: 'parent2',
        primeEstimatedWeight: 2000,
        status: shipmentStatuses.APPROVED,
        diversion: true,
        pickupAddress: { city: 'City1' },
        destinationAddress: { city: 'City2' },
      },
      {
        id: 'child2',
        primeEstimatedWeight: 1500,
        status: shipmentStatuses.APPROVED,
        diversion: true,
        pickupAddress: { city: 'City2' },
        destinationAddress: { city: 'City3' },
      },
    ];

    const { result } = renderHook(() => useCalculatedEstimatedWeight(mtoShipments));

    expect(result.current).toBe(2500 + 1500);
  });
  it('calculates the lowest weight of shipments with 1 parent diversion', () => {
    const mtoShipments = [
      {
        id: 'parent1',
        primeEstimatedWeight: 3000,
        status: shipmentStatuses.APPROVED,
        diversion: true,
        pickupAddress: { city: 'CityA' },
        destinationAddress: { city: 'CityB' },
      },
      {
        id: 'child1',
        primeEstimatedWeight: 2500,
        status: shipmentStatuses.APPROVED,
        diversion: true,
        pickupAddress: { city: 'CityB' },
        destinationAddress: { city: 'CityC' },
      },
    ];

    const { result } = renderHook(() => useCalculatedEstimatedWeight(mtoShipments));

    expect(result.current).toBe(2500);
  });
});
