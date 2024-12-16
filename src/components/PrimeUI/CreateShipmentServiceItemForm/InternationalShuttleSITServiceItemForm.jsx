import * as Yup from 'yup';
import { Formik } from 'formik';
import { Button } from '@trussworks/react-uswds';
import React from 'react';
import PropTypes from 'prop-types';

import TextField from 'components/form/fields/TextField/TextField';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { Form } from 'components/form/Form';
import { ShipmentShape } from 'types/shipment';
import { DropdownInput } from 'components/form/fields';
import { shuttleServiceItemCodeOptions, createServiceItemModelTypes } from 'constants/prime';

const internationalShuttleSITValidationSchema = Yup.object().shape({
  reServiceCode: Yup.string().required('Required'),
  reason: Yup.string().required('Required'),
});

const InternationalShuttleSITServiceItemForm = ({ shipment, submission }) => {
  const initialValues = {
    moveTaskOrderID: shipment.moveTaskOrderID,
    mtoShipmentID: shipment.id,
    modelType: createServiceItemModelTypes.MTOServiceItemInternationalShuttle,
    reason: '',
    estimatedWeight: null,
    actualWeight: null,
  };

  const onSubmit = (values) => {
    const { estimatedWeight, actualWeight, ...otherFields } = values;
    const body = {
      estimatedWeight: Number.parseInt(estimatedWeight, 10),
      actualWeight: Number.parseInt(actualWeight, 10),
      ...otherFields,
    };
    submission({ body });
  };

  return (
    <Formik
      initialValues={initialValues}
      validationSchema={internationalShuttleSITValidationSchema}
      onSubmit={onSubmit}
    >
      <Form data-testid="internationalShuttleSITServiceItemForm">
        <DropdownInput
          label="Service item code"
          name="reServiceCode"
          id="reServiceCode"
          required
          options={shuttleServiceItemCodeOptions}
        />
        <TextField name="reason" id="reason" label="Reason" />
        <MaskedTextField
          data-testid="estimatedWeightInput"
          name="estimatedWeight"
          label="Estimated weight (lbs)"
          labelHint="Optional"
          id="estimatedWeightInput"
          mask={Number}
          scale={0}
          thousandsSeparator=","
          lazy={false}
        />
        <MaskedTextField
          data-testid="actualWeightInput"
          name="actualWeight"
          label="Actual weight (lbs)"
          labelHint="Optional"
          id="actualWeightInput"
          mask={Number}
          scale={0}
          thousandsSeparator=","
          lazy={false}
        />
        <Button type="submit">Create service item</Button>
      </Form>
    </Formik>
  );
};

InternationalShuttleSITServiceItemForm.propTypes = {
  shipment: ShipmentShape.isRequired,
  submission: PropTypes.func.isRequired,
};

export default InternationalShuttleSITServiceItemForm;
