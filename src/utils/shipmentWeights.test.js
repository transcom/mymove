import {
  calculateNetWeightForProGearWeightTicket,
  calculateWeightTicketWeightDifference,
  calculateNonPPMShipmentNetWeight,
  calculatePPMShipmentNetWeight,
  calculateShipmentNetWeight,
  calculateTotalNetWeightForProGearWeightTickets,
  getTotalNetWeightForWeightTickets,
  getShipmentEstimatedWeight,
  shipmentIsOverweight,
  getWeightTicketNetWeight,
} from './shipmentWeights';
import {
  createCompleteProGearWeightTicket,
  createRejectedProGearWeightTicket,
} from './test/factories/proGearWeightTicket';
import { createCompleteWeightTicket } from './test/factories/weightTicket';

describe('shipmentWeights utils', () => {
  describe('shipmentIsOverweight', () => {
    it('returns true when the shipment weight is over 110% of the estimated weight', () => {
      expect(shipmentIsOverweight(100, 111)).toEqual(true);
    });

    it('returns false when shipment weight is less than  110% of the estimated weight', () => {
      expect(shipmentIsOverweight(100, 101)).toEqual(false);
    });

    it('returns false when estimated weight is undefined', () => {
      expect(shipmentIsOverweight(undefined, 100)).toEqual(false);
    });
  });
});

describe('calculateNetWeightForWeightTicket', () => {
  it.each([
    [0, 400, 400],
    [15000, 18000, 3000],
    [null, 1500, 0],
    [0, null, 0],
    [null, null, 0],
    [undefined, 1500, 0],
    [0, undefined, 0],
    [undefined, undefined, 0],
    ['not a number', 1500, 0],
    [0, 'not a number', 0],
    ['not a number', 'not a number', 0],
  ])(
    `calculates net weight properly | emptyWeight: %s | fullWeight: %s | expectedNetWeight: %s`,
    (emptyWeight, fullWeight, expectedNetWeight) => {
      const weightTicket = createCompleteWeightTicket(
        {},
        {
          emptyWeight,
          fullWeight,
        },
      );

      expect(calculateWeightTicketWeightDifference(weightTicket)).toEqual(expectedNetWeight);
    },
  );
});

describe('getTotalNetWeightForWeightTickets', () => {
  it.each([
    [[{ emptyWeight: 0, fullWeight: 400 }], 400],
    [
      [
        { emptyWeight: 0, fullWeight: 400 },
        { emptyWeight: 15000, fullWeight: 18000 },
      ],
      3400,
    ],
    [
      [
        { emptyWeight: null, fullWeight: 400 },
        { emptyWeight: 14000, fullWeight: 17000 },
      ],
      3000,
    ],
    [
      [
        { emptyWeight: 0, fullWeight: null },
        { emptyWeight: 14000, fullWeight: 19000 },
      ],
      5000,
    ],
    [
      [
        { emptyWeight: null, fullWeight: null },
        { emptyWeight: 14000, fullWeight: 18000 },
      ],
      4000,
    ],
    [
      [
        { emptyWeight: null, fullWeight: null },
        { emptyWeight: null, fullWeight: null },
      ],
      0,
    ],
    [
      [
        { emptyWeight: undefined, fullWeight: 400 },
        { emptyWeight: 14000, fullWeight: 17000 },
      ],
      3000,
    ],
    [
      [
        { emptyWeight: 0, fullWeight: undefined },
        { emptyWeight: 14000, fullWeight: 19000 },
      ],
      5000,
    ],
    [
      [
        { emptyWeight: undefined, fullWeight: undefined },
        { emptyWeight: 14000, fullWeight: 18000 },
      ],
      4000,
    ],
    [
      [
        { emptyWeight: undefined, fullWeight: undefined },
        { emptyWeight: undefined, fullWeight: undefined },
      ],
      0,
    ],
    [
      [
        { emptyWeight: 'not a number', fullWeight: 400 },
        { emptyWeight: 14000, fullWeight: 17000 },
      ],
      3000,
    ],
    [
      [
        { emptyWeight: 0, fullWeight: 'not a number' },
        { emptyWeight: 14000, fullWeight: 19000 },
      ],
      5000,
    ],
    [
      [
        { emptyWeight: 'not a number', fullWeight: 'not a number' },
        { emptyWeight: 14000, fullWeight: 18000 },
      ],
      4000,
    ],
    [
      [
        { emptyWeight: 'not a number', fullWeight: 'not a number' },
        { emptyWeight: 'not a number', fullWeight: 'not a number' },
      ],
      0,
    ],
    [
      [
        { emptyWeight: 'not a number', fullWeight: 'not a number' },
        { emptyWeight: 'not a number', fullWeight: 'not a number', adjustedNetWeight: 1000 },
      ],
      1000,
    ],
    [[], 0],
  ])(`calculates total net weight properly`, (weightTicketsFields, expectedNetWeight) => {
    const weightTickets = [];

    weightTicketsFields.forEach((fieldOverrides) => {
      weightTickets.push(createCompleteWeightTicket({}, fieldOverrides));
    });

    expect(getTotalNetWeightForWeightTickets(weightTickets)).toEqual(expectedNetWeight);
  });
});

describe('getWeightTicketNetWeight', () => {
  it('returns the adjusted net weight if present', () => {
    const weightTicket = { emptyWeight: 4, fullWeight: 10, adjustedNetWeight: 1000 };
    expect(getWeightTicketNetWeight(weightTicket)).toEqual(1000);
  });
  it('returns the calculated weight difference if the net weight has not been adjusted', () => {
    const weightTicket = { emptyWeight: 4, fullWeight: 10, adjustedNetWeight: null };
    expect(getWeightTicketNetWeight(weightTicket)).toEqual(6);
  });
});

describe('calculateNetWeightForProGearWeightTicket', () => {
  it.each([
    [0, 0],
    [15000, 15000],
    [null, 0],
    [undefined, 0],
    ['not a number', 0],
  ])(
    `calculates net weight properly | emptyWeight: %s | fullWeight: %s | constructedWeight: %s | expectedNetWeight: %s`,
    (weight, expectedNetWeight) => {
      const proGearWeightTicket = createCompleteProGearWeightTicket(
        {},
        {
          weight,
        },
      );

      expect(calculateNetWeightForProGearWeightTicket(proGearWeightTicket)).toEqual(expectedNetWeight);
    },
  );

  it('rejected weight ticket net weight is zero', () => {
    const rejectedProGearWeightTicket = createRejectedProGearWeightTicket({}, { weight: 200 });
    // The weight of the ticket should be greater than zero for this test to be valid.
    expect(rejectedProGearWeightTicket.weight).toBeGreaterThan(0);
    expect(calculateNetWeightForProGearWeightTicket(rejectedProGearWeightTicket)).toEqual(0);
  });
});

describe('calculateTotalNetWeightForProGearWeightTickets', () => {
  it.each([
    [[{ weight: 0 }], 0],
    [[{ weight: 0 }, { weight: 15000 }], 15000],
    [[{ weight: null }], 0],
    [[{ weight: null }, { weight: 15000 }], 15000],
    [[{ weight: undefined }], 0],
    [[{ weight: undefined }, { weight: 15000 }], 15000],
    [[{ weight: 'not a number' }], 0],
    [[{ weight: 'not a number' }, { weight: 15000 }], 15000],
    [[], 0],
  ])(`calculates total net weight properly`, (proGearWeightTicketsFields, expectedNetWeight) => {
    const proGearWeightTickets = [];

    proGearWeightTicketsFields.forEach((fieldOverrides) => {
      proGearWeightTickets.push(createCompleteProGearWeightTicket({}, fieldOverrides));
    });

    expect(calculateTotalNetWeightForProGearWeightTickets(proGearWeightTickets)).toEqual(expectedNetWeight);
  });
});

describe('calculateTotalNetWeightForProGearWeightTickets with a rejected weight ticket', () => {
  it('rejected weight ticket is not included in total net weight', () => {
    const approvedWeight = 350;
    const approvedProGearWeightTicket = createCompleteProGearWeightTicket({}, { weight: approvedWeight });
    const rejectedProGearWeightTicket = createRejectedProGearWeightTicket({}, { weight: 200 });
    // The weight of each ticket should be greater than zero for this test to be valid.
    expect(approvedProGearWeightTicket.weight).toBeGreaterThan(0);
    expect(rejectedProGearWeightTicket.weight).toBeGreaterThan(0);
    expect(
      calculateTotalNetWeightForProGearWeightTickets([approvedProGearWeightTicket, rejectedProGearWeightTicket]),
    ).toEqual(approvedWeight);
  });

  it('returns 0 if not tickets are passed in', () => {
    expect(calculateTotalNetWeightForProGearWeightTickets(null)).toEqual(0);
  });
});

describe('Calculating shipment net weights', () => {
  const ppmShipments = [
    {
      ppmShipment: {
        weightTickets: [{ emptyWeight: 14000, fullWeight: 19000 }],
      },
    },
    {
      ppmShipment: {
        weightTickets: [
          { emptyWeight: 14000, fullWeight: 19000 },
          { emptyWeight: 12000, fullWeight: 18000 },
        ],
      },
    },
    {
      ppmShipment: {
        weightTickets: [
          { emptyWeight: 14000, fullWeight: 19000 },
          { emptyWeight: 12000, fullWeight: 18000 },
          { emptyWeight: 10000, fullWeight: 20000 },
        ],
      },
    },
  ];

  const hhgShipments = [
    {
      primeActualWeight: 10,
      reweigh: {
        weight: 5,
      },
    },
    {
      primeActualWeight: 2000,
      reweigh: {
        weight: 300,
      },
    },
    {
      primeActualWeight: 100,
    },
    {
      primeActualWeight: 1000,
      reweigh: {
        weight: 200,
      },
    },
    {
      primeActualWeight: 400,
      reweigh: {
        weight: 3000,
      },
    },
  ];

  it('calculates the net weight of a ppm shipment properly', () => {
    expect(calculatePPMShipmentNetWeight(ppmShipments[0])).toEqual(5000);
  });

  it('calculates the net weight of a non-ppm shipment properly', () => {
    expect(calculateNonPPMShipmentNetWeight(hhgShipments[0])).toEqual(5);
  });

  it('calculates the sum net weight of a move with varied shipment types', () => {
    const netWeightOfPPMShipments = 37000;
    const netWeightOfNonPPMShipments = 1005;

    const totalMoveWeight = [...ppmShipments, ...hhgShipments]
      .map((s) => calculateShipmentNetWeight(s))
      .reduce((accumulator, current) => accumulator + current, 0);
    expect(totalMoveWeight).toEqual(netWeightOfPPMShipments + netWeightOfNonPPMShipments);
  });
});

describe('Calculating shipment estimated weights', () => {
  const ppmShipments = [
    {
      shipmentType: 'PPM',
      ppmShipment: {
        estimatedWeight: 5000,
      },
    },
    {
      shipmentType: 'PPM',
      ppmShipment: {
        estimatedWeight: 11000,
      },
    },
    {
      shipmentType: 'PPM',
      ppmShipment: {
        weightTickets: [
          { emptyWeight: 14000, fullWeight: 19000 },
          { emptyWeight: 12000, fullWeight: 18000 },
          { emptyWeight: 10000, fullWeight: 20000 },
        ],
      },
    },
  ];

  const hhgShipments = [
    {
      shipmentType: 'HHG',
      primeEstimatedWeight: 10,
    },
    {
      shipmentType: 'HHG',
      primeEstimatedWeight: 2000,
    },
    {
      shipmentType: 'HHG',
      primeEstimatedWeight: 100,
    },
    {
      shipmentType: 'HHG',
      primeEstimatedWeight: 1000,
    },
    {
      shipmentType: 'HHG',
      reweigh: {
        weight: 3000,
      },
    },
  ];

  it('gets the estimated weight of a ppm shipment properly', () => {
    expect(getShipmentEstimatedWeight(ppmShipments[0])).toEqual(5000);
  });

  it('gets the estimated weight of a non-ppm shipment properly', () => {
    expect(getShipmentEstimatedWeight(hhgShipments[0])).toEqual(10);
  });

  it('calculates the sum net weight of a move with varied shipment types', () => {
    const estimatedWeightOfPPMShipments = 16000;
    const estimatedWeightOfNonPPMShipments = 3110;

    const totalEstimatedWeight = [...ppmShipments, ...hhgShipments]
      .map((s) => getShipmentEstimatedWeight(s))
      .reduce((accumulator, current) => accumulator + current, 0);
    expect(totalEstimatedWeight).toEqual(estimatedWeightOfPPMShipments + estimatedWeightOfNonPPMShipments);
  });
});
