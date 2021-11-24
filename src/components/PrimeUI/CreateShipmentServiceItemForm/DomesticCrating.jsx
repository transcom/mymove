import * as Yup from 'yup';
import { Field, Formik } from 'formik';
import { Button, Label, Textarea } from '@trussworks/react-uswds';
import React from 'react';
import PropTypes from 'prop-types';

import { Form } from 'components/form/Form';
import TextField from 'components/form/fields/TextField';
import { DropdownInput } from 'components/form/fields/DropdownInput';
import { ShipmentShape } from 'types/shipment';
import { domesticCratingServiceItemCodeOptions, createServiceItemModelTypes } from 'constants/prime';
import MaskedTextField from 'components/form/fields/MaskedTextField';

const domesticShippingValidationSchema = Yup.object().shape({});

const DomesticShippingServiceItemForm = ({ shipment, submission }) => {
  const initialValues = {
    moveTaskOrderID: shipment.moveTaskOrderID,
    mtoShipmentID: shipment.id,
    modelType: createServiceItemModelTypes.MTOServiceItemDomesticCrating,
    itemLength: '0',
    itemWidth: '0',
    itemHeight: '0',
    crateLength: '0',
    crateWidth: '0',
    crateHeight: '0',
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
      <Form>
        <DropdownInput
          label="Service Item Code"
          name="reServiceCode"
          id="reServiceCode"
          required
          options={domesticCratingServiceItemCodeOptions}
        />
        <MaskedTextField
          data-testid="itemLength"
          defaultValue="0"
          name="itemLength"
          label="Item length (ft)"
          id="itemLength"
          mask={Number}
          scale={0} // digits after point, 0 for integers
          thousandsSeparator=","
          lazy={false} // immediate masking evaluation
        />
        <MaskedTextField
          data-testid="itemWidth"
          defaultValue="0"
          name="itemWidth"
          label="Item width (ft)"
          id="itemWidth"
          mask={Number}
          scale={0} // digits after point, 0 for integers
          thousandsSeparator=","
          lazy={false} // immediate masking evaluation
        />
        <MaskedTextField
          data-testid="itemHeight"
          defaultValue="0"
          name="itemHeight"
          label="Item height (ft)"
          id="itemHeight"
          mask={Number}
          scale={0} // digits after point, 0 for integers
          thousandsSeparator=","
          lazy={false} // immediate masking evaluation
        />
        <MaskedTextField
          data-testid="crateLength"
          defaultValue="0"
          name="crateLength"
          label="Crate length (ft)"
          id="crateLength"
          mask={Number}
          scale={0} // digits after point, 0 for integers
          thousandsSeparator=","
          lazy={false} // immediate masking evaluation
        />
        <MaskedTextField
          data-testid="crateWidth"
          defaultValue="0"
          name="crateWidth"
          label="Crate width (ft)"
          id="crateWidth"
          mask={Number}
          scale={0} // digits after point, 0 for integers
          thousandsSeparator=","
          lazy={false} // immediate masking evaluation
        />
        <MaskedTextField
          data-testid="crateHeight"
          defaultValue="0"
          name="crateHeight"
          label="Crate height (ft)"
          id="crateHeight"
          mask={Number}
          scale={0} // digits after point, 0 for integers
          thousandsSeparator=","
          lazy={false} // immediate masking evaluation
        />
        <Label htmlFor="description">Description</Label>
        <Field as={Textarea} data-testid="description" name="description" id="description" />
        <TextField name="reason" id="reason" label="Reason" />
        <Button type="submit">Create service item</Button>
      </Form>
    </Formik>
  );
};

DomesticShippingServiceItemForm.propTypes = {
  shipment: ShipmentShape.isRequired,
  submission: PropTypes.func.isRequired,
};

export default DomesticShippingServiceItemForm;
