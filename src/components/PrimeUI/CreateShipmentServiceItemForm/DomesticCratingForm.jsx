import * as Yup from 'yup';
import { Formik } from 'formik';
import { Button } from '@trussworks/react-uswds';
import React from 'react';
import PropTypes from 'prop-types';

import { Form } from 'components/form/Form';
import TextField from 'components/form/fields/TextField/TextField';
import { DropdownInput } from 'components/form/fields/DropdownInput';
import { ShipmentShape } from 'types/shipment';
import { domesticCratingServiceItemCodeOptions, createServiceItemModelTypes } from 'constants/prime';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';

const domesticShippingValidationSchema = Yup.object().shape({
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

const DomesticCratingForm = ({ shipment, submission }) => {
  const initialValues = {
    moveTaskOrderID: shipment.moveTaskOrderID,
    mtoShipmentID: shipment.id,
    modelType: createServiceItemModelTypes.MTOServiceItemDomesticCrating,
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
    <Formik initialValues={initialValues} validationSchema={domesticShippingValidationSchema} onSubmit={onSubmit}>
      <Form data-testid="domesticCratingForm">
        <DropdownInput
          label="Service item code"
          name="reServiceCode"
          id="reServiceCode"
          required
          options={domesticCratingServiceItemCodeOptions}
        />
        <MaskedTextField
          data-testid="itemLength"
          name="itemLength"
          label="Item length (ft)"
          id="itemLength"
          mask={Number}
          scale={0}
          thousandsSeparator=","
          lazy={false}
        />
        <MaskedTextField
          data-testid="itemWidth"
          name="itemWidth"
          label="Item width (ft)"
          id="itemWidth"
          mask={Number}
          scale={0}
          thousandsSeparator=","
          lazy={false}
        />
        <MaskedTextField
          data-testid="itemHeight"
          name="itemHeight"
          label="Item height (ft)"
          id="itemHeight"
          mask={Number}
          scale={0}
          thousandsSeparator=","
          lazy={false}
        />
        <MaskedTextField
          data-testid="crateLength"
          name="crateLength"
          label="Crate length (ft)"
          id="crateLength"
          mask={Number}
          scale={0}
          thousandsSeparator=","
          lazy={false}
        />
        <MaskedTextField
          data-testid="crateWidth"
          name="crateWidth"
          label="Crate width (ft)"
          id="crateWidth"
          mask={Number}
          scale={0}
          thousandsSeparator=","
          lazy={false}
        />
        <MaskedTextField
          data-testid="crateHeight"
          name="crateHeight"
          label="Crate height (ft)"
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
    </Formik>
  );
};

DomesticCratingForm.propTypes = {
  shipment: ShipmentShape.isRequired,
  submission: PropTypes.func.isRequired,
};

export default DomesticCratingForm;
