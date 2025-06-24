import React from 'react';
import PropTypes from 'prop-types';
import { Formik } from 'formik';
import { Button } from '@trussworks/react-uswds';
import * as Yup from 'yup';

import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import Modal, { ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';
import TextField from 'components/form/fields/TextField/TextField';
import { Form } from 'components/form';
import { requiredAsteriskMessage } from 'components/form/RequiredAsterisk';

const diversionReasonSchema = Yup.object().shape({
  diversionReason: Yup.string().required('Required'),
});

const RequestShipmentDiversionModal = ({ onClose, onSubmit, shipmentInfo }) => {
  let validDate = false;
  const today = new Date();
  const pickupDate = new Date(shipmentInfo.actualPickupDate);
  if (today >= pickupDate) {
    validDate = true;
  }

  return (
    <div>
      <Overlay />
      <ModalContainer>
        <Modal>
          <ModalClose handleClick={() => onClose()} />
          <ModalTitle>
            <h3>Request Shipment Diversion for #{shipmentInfo.shipmentLocator}</h3>
          </ModalTitle>
          <p>
            Movers will be notified that a diversion has been requested on this shipment. They will confirm or deny this
            request.
          </p>
          <Formik
            initialValues={{ diversionReason: '' }}
            validationSchema={diversionReasonSchema}
            onSubmit={(values) => {
              onSubmit(shipmentInfo.id, shipmentInfo.eTag, shipmentInfo.shipmentLocator, values.diversionReason);
            }}
          >
            {({ handleChange, values, isValid, dirty }) => {
              return (
                <Form aria-label="diversion reason">
                  {requiredAsteriskMessage}
                  <TextField
                    id="diversionReason"
                    name="diversionReason"
                    label="Reason for Diversion"
                    type="text"
                    value={values.diversionReason}
                    onChange={handleChange}
                    showRequiredAsterisk
                    required
                  />
                  <ModalActions>
                    <Button secondary type="reset" onClick={() => onClose()} data-testid="modalBackButton">
                      Back
                    </Button>
                    <Button type="submit" disabled={!isValid || !dirty || !validDate} data-testid="modalSubmitButton">
                      Request Diversion
                    </Button>
                  </ModalActions>
                </Form>
              );
            }}
          </Formik>
        </Modal>
      </ModalContainer>
    </div>
  );
};

RequestShipmentDiversionModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
  shipmentInfo: PropTypes.shape({
    shipmentID: PropTypes.string.isRequired,
    ifMatchEtag: PropTypes.string.isRequired,
    moveTaskOrderID: PropTypes.string.isRequired,
  }).isRequired,
};

RequestShipmentDiversionModal.displayName = 'RequestShipmentDiversionModal';

export default RequestShipmentDiversionModal;
