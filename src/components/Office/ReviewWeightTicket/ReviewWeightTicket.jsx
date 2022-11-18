import React, { useState } from 'react';
import { func, object } from 'prop-types';
import { Field, Formik } from 'formik';
import classnames from 'classnames';
import { Form, Radio } from '@trussworks/react-uswds';
// import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './ReviewWeightTicket.module.scss';

import ReviewItems from 'components/Customer/PPM/Closeout/ReviewItems/ReviewItems';
// import DocumentViewerSidebar from 'pages/Office/DocumentViewerSidebar/DocumentViewerSidebar';
import { formatAboutYourPPMItem } from 'utils/ppmCloseout';
// import officeShapes from 'types/officeShapes';
import SectionWrapper from 'components/Customer/SectionWrapper';
import Fieldset from 'shared/Fieldset';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import ApproveReject from 'components/form/ApproveReject/ApproveReject';

export default function ReviewWeightTicket({ mtoShipment, onClose }) {
  const aboutYourPPM = formatAboutYourPPMItem(mtoShipment?.ppmShipment);
  const [canEditRejection, setCanEditRejection] = useState(!rejectionReason);

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <div>
          <h1>PPM 1</h1>
          <ReviewItems contents={aboutYourPPM} />
        </div>
        <hr />
        <h2>Trip 1</h2>
        <form>
          <h3>Vehicle description</h3>
        </form>
        <Formik>
          {({ isValid, isSubmitting, handleReset, handleChange, submitForm, values }) => {
            const handleApprovalChange = (event) => {
              handleChange(event);
              submitForm().then(() => {
                setCanEditRejection(true);
              });
            };

            const handleRejectChange = (event) => {
              handleChange(event);
              submitForm().then(() => {
                setCanEditRejection(false);
              });
            };

            const handleRejectCancel = (event) => {
              if (initialValues.rejectionReason) {
                setCanEditRejection(false);
              }

              handleReset(event);
            };
            return (
              <div className={classnames(styles.AboutForm)}>
                <Form className={classnames(styles.W2Address)}>
                  <SectionWrapper>
                    <h2>Trip 1</h2>
                    <legend>Vehicle description</legend>
                    <div>Chevy</div>
                    <Fieldset>
                      <legend>Weight type</legend>
                      <Field
                        as={Radio}
                        id="weight-tickets"
                        label="Weight tickets"
                        name="weightTickets"
                        value="weightTicket"
                        checked
                      />
                      <Field
                        as={Radio}
                        id="constructed-weight"
                        label="Constructed weight"
                        name="constructedWeight"
                        value="constructedWeight"
                      />
                    </Fieldset>
                    <MaskedTextField
                      defaultValue="0"
                      name="emptyWeight"
                      label="Empty weight"
                      id="emptyWeight"
                      mask={Number}
                      scale={0} // digits after point, 0 for integers
                      signed={false} // disallow negative
                      thousandsSeparator=","
                      lazy={false} // immediate masking evaluation
                      suffix="lbs"
                    />
                    <MaskedTextField
                      defaultValue="0"
                      name="fullWeight"
                      label="Full weight"
                      id="fullWeight"
                      mask={Number}
                      scale={0} // digits after point, 0 for integers
                      signed={false} // disallow negative
                      thousandsSeparator=","
                      lazy={false} // immediate masking evaluation
                      suffix="lbs"
                    />
                    <legend>Net weight</legend>
                    <div>4,565 lbs</div>
                    <Fieldset>
                      <legend>Did they use a trailer they owned</legend>
                      <Field as={Radio} id="yes" label="Yes" name="yesOwnedTrailer" value="true" checked />
                      <Field as={Radio} id="no" label="No" name="noOwnTrailer" value="false" />
                    </Fieldset>
                    <ApproveReject
                      id={mtoShipment?.ppmShipment?.id}
                      currentStatus={values.status}
                      rejectionReason={values.rejectionReason}
                      requestComplete
                      approvedStatus={APPROVED}
                      deniedStatus={DENIED}
                      canEditRejection={canEditRejection}
                      setCanEditRejection={setCanEditRejection}
                      handleApprovalChange={handleApprovalChange}
                      handleRejectChange={handleRejectChange}
                      handleRejectCancel={handleRejectCancel}
                      handleChange={handleChange}
                      handleFormReset={handleFormReset}
                    />
                  </SectionWrapper>
                </Form>
              </div>
            );
          }}
        </Formik>
      </header>
    </div>
  );
}

ReviewWeightTicket.propTypes = {
  mtoShipment: object.isRequired,
  onClose: func.isRequired,
};
