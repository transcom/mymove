import * as Yup from 'yup';
import { Formik } from 'formik';
import { Button } from '@trussworks/react-uswds';
import React from 'react';
import PropTypes from 'prop-types';

import { Form } from 'components/form/Form';
import TextField from 'components/form/fields/TextField';
import { ShipmentShape } from 'types/shipment';
import { DropdownInput } from 'components/form/fields';
import { shuttleServiceItemCodeOptions, createServiceItemModelTypes } from 'constants/prime';
import MaskedTextField from 'components/form/fields/MaskedTextField';

const shuttleSITValidationSchema = Yup.object().shape({});

const ShuttleSITServiceItemForm = ({ shipment, submission }) => {
  const initialValues = {
    moveTaskOrderID: shipment.moveTaskOrderID,
    mtoShipmentID: shipment.id,
    modelType: createServiceItemModelTypes.MTOServiceItemShuttle,
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
    <Formik initialValues={initialValues} validationSchema={shuttleSITValidationSchema} onSubmit={onSubmit}>
      <Form>
        <DropdownInput
          label="Service Item Code"
          name="reServiceCode"
          id="reServiceCode"
          required
          options={shuttleServiceItemCodeOptions}
        />
        <TextField name="reason" id="reason" label="Reason" />
        <MaskedTextField
          data-testid="estimatedWeightInput"
          defaultValue="0"
          name="estimatedWeight"
          label="Estimated weight (lbs)"
          id="estimatedWeightInput"
          mask={Number}
          scale={0} // digits after point, 0 for integers
          thousandsSeparator=","
          lazy={false} // immediate masking evaluation
        />
        <MaskedTextField
          data-testid="actualWeightInput"
          defaultValue="0"
          name="actualWeight"
          label="Actual weight (lbs)"
          id="actualWeightInput"
          mask={Number}
          scale={0} // digits after point, 0 for integers
          thousandsSeparator=","
          lazy={false} // immediate masking evaluation
        />
        <Button type="submit">Create service item</Button>
      </Form>
    </Formik>
  );
};

ShuttleSITServiceItemForm.propTypes = {
  shipment: ShipmentShape.isRequired,
  submission: PropTypes.func.isRequired,
};

export default ShuttleSITServiceItemForm;
