import React from 'react';
import { render, screen } from '@testing-library/react';

import styles from './ReviewShipmentWeightsTable.module.scss';
import {
  addShipmentNumbersToTableData,
  determineTableRowClassname,
  ShipmentTypeCell,
  estimatedWeightDisplayHelper,
  DASH,
  actualWeightDisplayHelper,
} from './helpers';

describe('addShipmentNumbersToTableData', () => {
  it.each([
    [[{ entitlements: { progearWeight: 2000 } }], false, [{ entitlements: { progearWeight: 2000 } }]],
    [[{ shipmentType: 'HHG' }], true, [{ shipmentType: 'HHG', showNumber: false, shipmentNumber: 1 }]],
    [
      [{ shipmentType: 'HHG' }, { shipmentType: 'HHG' }],
      true,
      [
        { shipmentType: 'HHG', showNumber: true, shipmentNumber: 1 },
        { shipmentType: 'HHG', showNumber: true, shipmentNumber: 2 },
      ],
    ],
  ])('correctly modifies the tableData', (tableData, determineShipmentNumbers, expectedResult) => {
    expect(addShipmentNumbersToTableData(tableData, determineShipmentNumbers)).toStrictEqual(expectedResult);
  });
});

describe('determineTableRowClassname', () => {
  it.each([
    ['HHG_OUTOF_NTS_DOMESTIC', styles[`review-shipment-weights-table-row-NTS-release`]],
    ['HHG_INTO_NTS', styles[`review-shipment-weights-table-row-NTS`]],
    ['PPM', styles[`review-shipment-weights-table-row-PPM`]],
    ['HHG', styles[`review-shipment-weights-table-row-HHG`]],
    ['NOT_AN_OPTION', ''],
  ])('returns the correct classname', (shipmentType, expectedClassname) => {
    expect(determineTableRowClassname(shipmentType)).toBe(expectedClassname);
  });
});

describe('shipmentTypeCellDisplayHelper', () => {
  it.each([
    [{ shipmentType: 'PPM', showNumber: false }, 'PPM'],
    [{ shipmentType: 'HHG_OUTOF_NTS_DOMESTIC', showNumber: true, shipmentNumber: 123 }, 'NTS-release 123'],
    [{ shipmentType: 'HHG', showNumber: true, shipmentNumber: 8 }, 'HHG 8'],
  ])('renders the correct Shipment Type Cell', (row, expectedResult) => {
    render(<ShipmentTypeCell row={row} />);
    expect(screen.getByText(expectedResult)).toBeInTheDocument();
  });
});

describe('estimatedWeightDisplayHelper', () => {
  it.each([
    [{ shipmentType: 'HHG_OUTOF_NTS_DOMESTIC' }, 'N/A'],
    [{ shipmentType: 'HHG_INTO_NTS', ntsRecordedWeight: 1234, primeEstimatedWeight: 9876 }, '1,234 lbs'],
    [{ shipmentType: 'HHG', ntsRecordedWeight: 1234, primeEstimatedWeight: 9876 }, '9,876 lbs'],
    [{ shipmentType: 'HHG', primeEstimatedWeight: 0 }, DASH],
  ])('renders the correct Shipment Type Cell', (row, expectedResult) => {
    expect(estimatedWeightDisplayHelper(row)).toBe(expectedResult);
  });
});

describe('actualWeightDisplayHelper', () => {
  it.each([
    [{ shipmentType: 'HHG_OUTOF_NTS_DOMESTIC' }, DASH],
    [{ primeActualWeight: 1234 }, '1,234 lbs'],
    [{ reweigh: { weight: 9876 } }, '9,876 lbs'],
    [{ primeActualWeight: 1234, reweigh: { weight: 9876 } }, '1,234 lbs'],
    [{ primeActualWeight: 9876, reweigh: { weight: 1234 } }, '1,234 lbs'],
  ])('renders the correct Shipment Type Cell', (row, expectedResult) => {
    expect(actualWeightDisplayHelper(row)).toBe(expectedResult);
  });
});
