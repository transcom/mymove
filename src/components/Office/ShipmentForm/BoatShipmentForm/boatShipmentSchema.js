import * as Yup from 'yup';

import {
  AdditionalAddressSchema,
  RequiredPlaceSchema,
  OptionalPlaceSchema,
} from 'components/Customer/MtoShipmentForm/validationSchemas';
import { toTotalInches } from 'utils/formatMtoShipment';

const currentYear = new Date().getFullYear();
const maxYear = currentYear + 2;

const boatShipmentSchema = () => {
  const formSchema = Yup.object()
    .shape({
      year: Yup.number().required('Required').min(1700, 'Invalid year').max(maxYear, 'Invalid year'),

      make: Yup.string().required('Required'),

      model: Yup.string().required('Required'),

      lengthFeet: Yup.number()
        .min(0)
        .nullable()
        .when('lengthInches', {
          is: (lengthInches) => !lengthInches,
          then: (schema) => schema.required('Required'),
          otherwise: (schema) => schema.notRequired(),
        }),

      lengthInches: Yup.number().min(0).nullable(),

      widthFeet: Yup.number()
        .min(0)
        .nullable()
        .when('widthInches', {
          is: (widthInches) => !widthInches,
          then: (schema) => schema.required('Required'),
          otherwise: (schema) => schema.notRequired(),
        }),

      widthInches: Yup.number().min(0).nullable(),

      heightFeet: Yup.number()
        .min(0)
        .nullable()
        .when('heightInches', {
          is: (heightInches) => !heightInches,
          then: (schema) => schema.required('Required'),
          otherwise: (schema) => schema.notRequired(),
        }),

      heightInches: Yup.number().min(0).nullable(),

      hasTrailer: Yup.boolean().required('Required'),

      isRoadworthy: Yup.boolean().when('hasTrailer', {
        is: true,
        then: (schema) => schema.required('Required'),
        otherwise: (schema) => schema.notRequired(),
      }),

      type: Yup.string().required('Required'),

      pickup: RequiredPlaceSchema,
      delivery: OptionalPlaceSchema,
      secondaryPickup: AdditionalAddressSchema,
      secondaryDelivery: AdditionalAddressSchema,
      tertiaryPickup: AdditionalAddressSchema,
      tertiaryDelivery: AdditionalAddressSchema,
      counselorRemarks: Yup.string(),
      customerRemarks: Yup.string(),
    })
    .test('dimension-check', 'Dimensions requirements.', function dimensionTest(values) {
      const { lengthFeet, lengthInches, widthFeet, widthInches, heightFeet, heightInches } = values;
      const hasLength = lengthFeet !== undefined || lengthInches !== undefined;
      const hasWidth = widthFeet !== undefined || widthInches !== undefined;
      const hasHeight = heightFeet !== undefined || heightInches !== undefined;

      if (hasLength && hasWidth && hasHeight) {
        const lengthInInches = toTotalInches(lengthFeet, lengthInches);
        const widthInInches = toTotalInches(widthFeet, widthInches);
        const heightInInches = toTotalInches(heightFeet, heightInches);

        if (lengthInInches <= 168 && widthInInches <= 82 && heightInInches <= 77) {
          const errors = [];
          errors.push(
            this.createError({
              path: 'lengthFeet',
              message: 'Dimensions do not meet the requirement.',
            }),
          );

          errors.push(
            this.createError({
              path: 'widthFeet',
              message: 'Dimensions do not meet the requirement.',
            }),
          );

          errors.push(
            this.createError({
              path: 'heightFeet',
              message: 'Dimensions do not meet the requirement.',
            }),
          );

          if (errors.length) {
            throw new Yup.ValidationError(errors);
          }
        }
      }
      return true;
    });
  return formSchema;
};

export default boatShipmentSchema;
