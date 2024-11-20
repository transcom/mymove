import * as Yup from 'yup';
import { Formik } from 'formik';
import { Button } from '@trussworks/react-uswds';
import React from 'react';
import PropTypes from 'prop-types';

import { Form } from 'components/form/Form';
import TextField from 'components/form/fields/TextField/TextField';
import { DropdownInput } from 'components/form/fields/DropdownInput';
import { ShipmentShape } from 'types/shipment';
import { internationalCratingServiceItemCodeOptions, createServiceItemModelTypes } from 'constants/prime';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { CheckboxField } from 'components/form/fields';
import { SERVICE_ITEM_CODES } from 'constants/serviceItems';

const internationalShippingValidationSchema = Yup.object().shape({
  reServiceCode: Yup.string().required('Required'),
  itemLength: Yup.string().required('Required'),
  itemWidth: Yup.string().required('Required'),
  itemHeight: Yup.string().required('Required'),
  crateLength: Yup.string().required('Required'),
  crateWidth: Yup.string().required('Required'),
  crateHeight: Yup.string().required('Required'),
  description: Yup.string().required('Required'),
  reason: Yup.string().required('Required'),
});

const InternationalCratingForm = ({ shipment, submission }) => {
  const initialValues = {
    moveTaskOrderID: shipment.moveTaskOrderID,
    mtoShipmentID: shipment.id,
    modelType: createServiceItemModelTypes.MTOServiceItemInternationalCrating,
    standaloneCrate: false,
    externalCrate: false,
    itemLength: '',
    itemWidth: '',
    itemHeight: '',
    crateLength: '',
    crateWidth: '',
    crateHeight: '',
    reason: '',
    description: '',
  };

  const onSubmit = (values) => {
    const { itemLength, itemWidth, itemHeight, crateLength, crateWidth, crateHeight, ...otherFields } = values;

    const body = {
      item: {
        length: Number.parseInt(itemLength, 10),
        width: Number.parseInt(itemWidth, 10),
        height: Number.parseInt(itemHeight, 10),
      },
      crate: {
        length: Number.parseInt(crateLength, 10),
        width: Number.parseInt(crateWidth, 10),
        height: Number.parseInt(crateHeight, 10),
      },
      ...otherFields,
    };
    submission({ body });
  };

  return (
    <Formik initialValues={initialValues} validationSchema={internationalShippingValidationSchema} onSubmit={onSubmit}>
      {({ values }) => {
        return (
          <Form data-testid="internationalCratingForm">
            <DropdownInput
              label="Service item code"
              name="reServiceCode"
              id="reServiceCode"
              required
              options={internationalCratingServiceItemCodeOptions}
            />
            {values.reServiceCode === SERVICE_ITEM_CODES.ICRT && (
              <>
                <CheckboxField id="standaloneCrate" name="standaloneCrate" label="Standalone Crate" />
                <CheckboxField id="externalCrate" name="externalCrate" label="External Crate" />
              </>
            )}
            <MaskedTextField
              data-testid="itemLength"
              name="itemLength"
              label="Item length (thousandths of an inch)"
              id="itemLength"
              mask={Number}
              scale={0}
              thousandsSeparator=","
              lazy={false}
            />
            <MaskedTextField
              data-testid="itemWidth"
              name="itemWidth"
              label="Item width (thousandths of an inch)"
              id="itemWidth"
              mask={Number}
              scale={0}
              thousandsSeparator=","
              lazy={false}
            />
            <MaskedTextField
              data-testid="itemHeight"
              name="itemHeight"
              label="Item height (thousandths of an inch)"
              id="itemHeight"
              mask={Number}
              scale={0}
              thousandsSeparator=","
              lazy={false}
            />
            <MaskedTextField
              data-testid="crateLength"
              name="crateLength"
              label="Crate length (thousandths of an inch)"
              id="crateLength"
              mask={Number}
              scale={0}
              thousandsSeparator=","
              lazy={false}
            />
            <MaskedTextField
              data-testid="crateWidth"
              name="crateWidth"
              label="Crate width (thousandths of an inch)"
              id="crateWidth"
              mask={Number}
              scale={0}
              thousandsSeparator=","
              lazy={false}
            />
            <MaskedTextField
              data-testid="crateHeight"
              name="crateHeight"
              label="Crate height (thousandths of an inch)"
              id="crateHeight"
              mask={Number}
              scale={0}
              thousandsSeparator=","
              lazy={false}
            />
            <TextField name="description" id="description" label="Description" />
            <TextField name="reason" id="reason" label="Reason" />
            <Button type="submit">Create service item</Button>
          </Form>
        );
      }}
    </Formik>
  );
};

InternationalCratingForm.propTypes = {
  shipment: ShipmentShape.isRequired,
  submission: PropTypes.func.isRequired,
};

export default InternationalCratingForm;
