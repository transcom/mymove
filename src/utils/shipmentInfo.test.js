import MOVE_STATUSES from 'constants/moves';
import determineShipmentInfo from 'utils/shipmentInfo';
import { SHIPMENT_OPTIONS } from 'shared/constants';

describe('determineShipmentInfo', () => {
  const fakeMove = {
    id: 'fakeMoveID',
    personally_procured_moves: [],
    status: MOVE_STATUSES.DRAFT,
  };

  it.each([
    [true, MOVE_STATUSES.DRAFT],
    [false, MOVE_STATUSES.SUBMITTED],
    [false, MOVE_STATUSES.APPROVALS_REQUESTED],
    [false, MOVE_STATUSES.APPROVED],
    [false, MOVE_STATUSES.CANCELED],
  ])('should set isHhgSelectable to %s if move is "%s"', (expectedValue, moveStatus) => {
    const move = { ...fakeMove, status: moveStatus };

    const info = determineShipmentInfo(move, []);

    expect(info.isHHGSelectable).toBe(expectedValue);
  });

  it.each([
    [true, []],
    [false, ['fakePPM']],
  ])('should set isPPMSelectable to %s if move PPM === %s', (expectedPPMSelectable, ppmList) => {
    const move = { ...fakeMove, personally_procured_moves: ppmList };

    const info = determineShipmentInfo(move, []);

    expect(info.isPPMSelectable).toBe(expectedPPMSelectable);
  });

  it.each([
    [1, [], []],
    [2, ['fakePPM'], []],
    [3, ['fakePPM', 'anotherPPM'], []],
    [3, ['fakePPM'], ['fakeMTO']],
    [5, ['fakePPM', 'anotherPPM'], ['fakeMTO', 'anotherMTO']],
    [4, ['fakePPM'], ['fakeMTO', 'anotherMTO']],
    [2, [], ['fakeMTO']],
    [3, [], ['fakeMTO', 'anotherMTO']],
  ])(
    'should return the correct new shipment number (%i) based on PPM (%s) and MTO (%s) shipments',
    (expectedNumber, ppmList, mtoShipments) => {
      const move = { ...fakeMove, personally_procured_moves: ppmList };

      const info = determineShipmentInfo(move, mtoShipments);

      expect(info.shipmentNumber).toBe(expectedNumber);
    },
  );

  it.each([
    [false, [], []],
    [true, ['fakePPM'], []],
    [true, ['fakePPM', 'anotherPPM'], []],
    [true, ['fakePPM'], ['fakeMTO']],
    [true, ['fakePPM'], ['fakeMTO', 'anotherMTO']],
    [true, [], ['fakeMTO']],
    [true, [], ['fakeMTO', 'anotherMTO']],
  ])(
    'should set hasShipment to (%i) based on PPM (%s) and MTO (%s) shipments',
    (expectedHasShipment, ppmList, mtoShipments) => {
      const move = { ...fakeMove, personally_procured_moves: ppmList };

      const info = determineShipmentInfo(move, mtoShipments);

      expect(info.hasShipment).toBe(expectedHasShipment);
    },
  );

  it.each([
    [true, MOVE_STATUSES.DRAFT, []],
    [true, MOVE_STATUSES.DRAFT, [{ shipmentType: SHIPMENT_OPTIONS.HHG }]],
    [true, MOVE_STATUSES.DRAFT, [{ shipmentType: SHIPMENT_OPTIONS.PPM }]],
    [false, MOVE_STATUSES.DRAFT, [{ shipmentType: SHIPMENT_OPTIONS.NTS }]],
    [false, MOVE_STATUSES.SUBMITTED, [{ shipmentType: SHIPMENT_OPTIONS.HHG }]],
    [false, MOVE_STATUSES.SUBMITTED, [{ shipmentType: SHIPMENT_OPTIONS.PPM }]],
    [true, MOVE_STATUSES.DRAFT, [{ shipmentType: SHIPMENT_OPTIONS.PPM }, { shipmentType: SHIPMENT_OPTIONS.HHG }]],
    [false, MOVE_STATUSES.DRAFT, [{ shipmentType: SHIPMENT_OPTIONS.PPM }, { shipmentType: SHIPMENT_OPTIONS.NTS }]],
    [false, MOVE_STATUSES.DRAFT, [{ shipmentType: SHIPMENT_OPTIONS.HHG }, { shipmentType: SHIPMENT_OPTIONS.NTS }]],
  ])(
    'sets isNTSSelectable to %s if move status is "%s" and MTO shipments === %s',
    (expectedNTSSelectable, moveStaus, mtoShipments) => {
      const move = { ...fakeMove, status: moveStaus };

      const info = determineShipmentInfo(move, mtoShipments);

      expect(info.isNTSSelectable).toBe(expectedNTSSelectable);
    },
  );

  it.each([
    [true, MOVE_STATUSES.DRAFT, []],
    [true, MOVE_STATUSES.DRAFT, [{ shipmentType: SHIPMENT_OPTIONS.HHG }]],
    [true, MOVE_STATUSES.DRAFT, [{ shipmentType: SHIPMENT_OPTIONS.PPM }]],
    [false, MOVE_STATUSES.DRAFT, [{ shipmentType: SHIPMENT_OPTIONS.NTSR }]],
    [false, MOVE_STATUSES.SUBMITTED, [{ shipmentType: SHIPMENT_OPTIONS.HHG }]],
    [false, MOVE_STATUSES.SUBMITTED, [{ shipmentType: SHIPMENT_OPTIONS.PPM }]],
    [true, MOVE_STATUSES.DRAFT, [{ shipmentType: SHIPMENT_OPTIONS.PPM }, { shipmentType: SHIPMENT_OPTIONS.HHG }]],
    [false, MOVE_STATUSES.DRAFT, [{ shipmentType: SHIPMENT_OPTIONS.PPM }, { shipmentType: SHIPMENT_OPTIONS.NTSR }]],
    [false, MOVE_STATUSES.DRAFT, [{ shipmentType: SHIPMENT_OPTIONS.HHG }, { shipmentType: SHIPMENT_OPTIONS.NTSR }]],
  ])(
    'sets isNTSRSelectable to %s if move status is "%s" and MTO shipments === %s',
    (expectedNTSRSelectable, moveStaus, mtoShipments) => {
      const move = { ...fakeMove, status: moveStaus };

      const info = determineShipmentInfo(move, mtoShipments);

      expect(info.isNTSRSelectable).toBe(expectedNTSRSelectable);
    },
  );
});
