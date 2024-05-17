import * as Yup from 'yup';

import { getFormattedMaxAdvancePercentage } from 'utils/incentives';
import { requiredAddressSchema } from 'utils/validation';
import { OptionalAddressSchema } from 'components/Customer/MtoShipmentForm/validationSchemas';
import { ADVANCE_STATUSES } from 'constants/ppms';

function closeoutOfficeSchema(showCloseoutOffice, isAdvancePage) {
  if (showCloseoutOffice && !isAdvancePage) {
    return Yup.object().shape({
      name: Yup.string().required('Required'),
    });
  }
  return Yup.object().notRequired();
}

const yupAdvancedPageSchemaGroup = (isOnAdvancedPage, theSchema, schemaWrapperFn) =>
  theSchema.when(['advanceRequested'], (_, schema) => {
    return isOnAdvancedPage ? schemaWrapperFn(schema) : schema;
  });

const applyNestedAdvanceValidator =
  ({ advanceAmountRequested, oldAdvanceStatus, hasRequestedAdvance }) =>
  (theSchema) =>
    theSchema
      .when(['advance', 'advanceStatus', 'advanceRequested'], {
        is: (advance, advanceStatus, advanceRequest) => {
          const amountsUnchanged = Number(advance) === advanceAmountRequested / 100;
          const statusUnchanged = advanceStatus === oldAdvanceStatus;
          const nonEditedInitialStatus = ADVANCE_STATUSES[advanceStatus] !== ADVANCE_STATUSES.EDITED;
          const requestUnchanged = advanceRequest === hasRequestedAdvance;
          const shouldShowError = amountsUnchanged && requestUnchanged && (statusUnchanged || !nonEditedInitialStatus);
          return !shouldShowError;
        },
        then: (schema) => {
          // console.log('logging for the advance request YES condition');
          return schema.required('Required');
        },
      })
      .when(['advanceRequested'], {
        is: (advanceRequest) => advanceRequest === false && advanceRequest !== hasRequestedAdvance,
        then: (schema) => {
          // console.log('logging for the advance request NO condition');
          return schema.required('Required');
        },
      });

const ppmShipmentSchema = ({
  estimatedIncentive = 0,
  weightAllotment = {},
  advanceAmountRequested = 0,
  hasRequestedAdvance,
  isAdvancePage,
  oldAdvanceStatus = '',
  showCloseoutOffice,
}) => {
  const estimatedWeightLimit = weightAllotment.totalWeightSelf || 0;
  const proGearWeightLimit = weightAllotment.proGearWeight || 0;
  const proGearSpouseWeightLimit = weightAllotment.proGearWeightSpouse || 0;

  const formSchema = Yup.object().shape({
    pickup: Yup.object().shape({
      address: requiredAddressSchema,
    }),
    destination: Yup.object().shape({
      address: requiredAddressSchema,
    }),
    secondaryPickup: Yup.object().shape({
      address: OptionalAddressSchema,
    }),
    secondaryDestination: Yup.object().shape({
      address: OptionalAddressSchema,
    }),

    expectedDepartureDate: Yup.date()
      .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
      .required('Required'),
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
        otherwise: (schema) =>
          schema
            .min(0, 'Enter a weight 0 lbs or greater')
            .max(proGearWeightLimit, `Enter a weight ${proGearWeightLimit.toLocaleString()} lbs or less`),
      }),
    spouseProGearWeight: Yup.number()
      .min(0, 'Enter a weight 0 lbs or greater')
      .max(proGearSpouseWeightLimit, `Enter a weight ${proGearSpouseWeightLimit.toLocaleString()} lbs or less`),
    closeoutOffice: closeoutOfficeSchema(showCloseoutOffice, isAdvancePage),

    advance: yupAdvancedPageSchemaGroup(isAdvancePage, Yup.number(), (theSchema) =>
      theSchema
        .max(
          (estimatedIncentive * 0.6) / 100,
          `Enter an amount that is less than or equal to the maximum advance (${getFormattedMaxAdvancePercentage()} of estimated incentive)`,
        )
        .min(1, 'Enter an amount $1 or more.')
        .when('advanceRequested', {
          is: true,
          then: (schema) => schema.required('Required'),
        }),
    ),
    counselorRemarks: yupAdvancedPageSchemaGroup(
      isAdvancePage,
      Yup.string(),
      applyNestedAdvanceValidator({ advanceAmountRequested, oldAdvanceStatus, hasRequestedAdvance }),
    ),
    // is: (advanceRequested) => {
    //   console.log(advanceRequested, hasRequestedAdvance);
    //   const unchangedRequestChoice = advanceRequested === hasRequestedAdvance;
    //   const shouldShowError = advanceRequested && unchangedRequestChoice;
    //   return isAdvancePage && !shouldShowError;
    // },
    // then: (schema) => {
    //   console.log(`requested schema valid ${schema.isValidSync({ recursive: true })}`);
    //   return schema.required('Required');
    // },
  });

  return formSchema;
};

export default ppmShipmentSchema;
