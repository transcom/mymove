import * as Yup from 'yup';

import { getFormattedMaxAdvancePercentage } from 'utils/incentives';
import { InvalidZIPTypeError, ZIP5_CODE_REGEX } from 'utils/validation';

const ppmShipmentSchema = ({
  estimatedIncentive = 0,
  weightAllotment = {},
  requestedAdvanceAmount = 0,
  isAdvanceRequested,
}) => {
  const estimatedWeightLimit = weightAllotment.totalWeightSelf || 0;
  const proGearWeightLimit = weightAllotment.proGearWeight || 0;
  const proGearSpouseWeightLimit = weightAllotment.proGearWeightSpouse || 0;

  const formSchema = Yup.object().shape({
    pickupPostalCode: Yup.string().matches(ZIP5_CODE_REGEX, InvalidZIPTypeError).required('Required'),
    secondPickupPostalCode: Yup.string().matches(ZIP5_CODE_REGEX, InvalidZIPTypeError),
    expectedDepartureDate: Yup.date()
      .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
      .required('Required'),

    destinationPostalCode: Yup.string().matches(ZIP5_CODE_REGEX, InvalidZIPTypeError).required('Required'),
    secondDestinationPostalCode: Yup.string().matches(ZIP5_CODE_REGEX, InvalidZIPTypeError),

    sitExpected: Yup.boolean().required('Required'),
    sitEstimatedWeight: Yup.number().when('sitExpected', {
      is: true,
      then: (schema) => schema.required('Required'),
    }),
    sitEstimatedEntryDate: Yup.date()
      .when('sitExpected', {
        is: true,
        then: (schema) =>
          schema.typeError('Enter a complete date in DD MMM YYYY format (day, month, year).').required('Required'),
      })
      .nullable(),
    sitEstimatedDepartureDate: Yup.date()
      .when('sitExpected', {
        is: true,
        then: (schema) =>
          schema.typeError('Enter a complete date in DD MMM YYYY format (day, month, year).').required('Required'),
      })
      .nullable(),

    estimatedWeight: Yup.number()
      .min(1, 'Enter a weight greater than 0 lbs')
      .max(estimatedWeightLimit, "Note: This weight exceeds the customer's weight allowance.")
      .required('Required'),
    hasProGear: Yup.boolean().required('Required'),
    proGearWeight: Yup.number()
      .min(0, 'Enter a weight 0 lbs or greater')
      .when(['hasProGear', 'spouseProGearWeight'], {
        is: (hasProGear, spouseProGearWeight) => hasProGear && !spouseProGearWeight,
        then: (schema) =>
          schema
            .required(
              `Enter weight in at least one pro-gear field. If the customer will not move pro-gear in this PPM, select No above.`,
            )
            .max(proGearWeightLimit, `Enter a weight ${proGearWeightLimit.toLocaleString()} lbs or less`),
        otherwise: Yup.number()
          .min(0, 'Enter a weight 0 lbs or greater')
          .max(proGearWeightLimit, `Enter a weight ${proGearWeightLimit.toLocaleString()} lbs or less`),
      }),
    spouseProGearWeight: Yup.number()
      .min(0, 'Enter a weight 0 lbs or greater')
      .max(proGearSpouseWeightLimit, `Enter a weight ${proGearSpouseWeightLimit.toLocaleString()} lbs or less`),

    advance: Yup.number()
      .max(
        (estimatedIncentive * 0.6) / 100,
        `Enter an amount that is less than or equal to the maximum advance (${getFormattedMaxAdvancePercentage()} of estimated incentive)`,
      )
      .when('advanceRequested', {
        is: true,
        then: (schema) => schema.required('Required'),
      }),

    counselorRemarks: Yup.string().when(['advance', 'advanceRequested'], {
      is: (advance, advanceRequested) =>
        Number(advance) !== requestedAdvanceAmount / 100 ||
        advanceRequested?.toString() !== isAdvanceRequested?.toString(),
      then: (schema) => schema.required('Required'),
    }),
  });

  return formSchema;
};

export default ppmShipmentSchema;
