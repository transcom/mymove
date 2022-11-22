import React, { useState } from 'react';
import { object, string } from 'prop-types';
import { Field, Formik } from 'formik';
import classnames from 'classnames';
import { Form, FormGroup, Label, Radio } from '@trussworks/react-uswds';

import styles from './ReviewWeightTicket.module.scss';

import Fieldset from 'shared/Fieldset';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import ApproveReject from 'components/form/ApproveReject/ApproveReject';
import formStyles from 'styles/form.module.scss';

export default function ReviewWeightTicket({ mtoShipment, tripNumber, ppmNumber }) {
  const [canEditRejection, setCanEditRejection] = useState(false);
  const initialValues = {
    weightType: '',
    emptyWeight: '',
    fullWeight: '',
    ownTrailer: '',
    status: '',
    rejectionReason: '',
  };
  return (
    <div className={styles.container}>
      <Formik initialValues={initialValues}>
        {({ handleReset, handleChange, submitForm, setValues, values }) => {
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

          const handleFormReset = () => {
            setValues({
              status: 'REQUESTED',
              rejectionReason: '',
            });
            submitForm().then(() => {
              setCanEditRejection(true);
            });
          };

          return (
            <div className="container--accent--ppm">
              <Form className={classnames(formStyles.form, styles.ReviewWeightTicket)}>
                <header className={styles.header}>
                  <div>
                    <h2>PPM {ppmNumber}</h2>
                    <section>
                      <div>
                        <Label className={styles.headerLabel}>Departure date</Label>
                        <span className={styles.light}>08-31-1991</span>
                      </div>
                      <div>
                        <Label className={styles.headerLabel}>Starting ZIP</Label>
                        <span className={styles.light}>90210</span>
                      </div>
                      <div>
                        <Label className={styles.headerLabel}>Ending ZIP</Label>
                        <span className={styles.light}>94611</span>
                      </div>
                      <div>
                        <Label className={styles.headerLabel}>Advance recieved</Label>
                        <span className={styles.light}>Yes, $560</span>
                      </div>
                    </section>
                  </div>
                  <hr />
                </header>
                <h3 className={styles.tripNumber}>Trip {tripNumber}</h3>
                <legend className={classnames('usa-label', styles.label)}>Vehicle description</legend>
                <div className={styles.displayValue}>Chevy</div>
                <FormGroup>
                  <Fieldset>
                    <legend className="usa-label">Weight type</legend>
                    <Field
                      as={Radio}
                      id="weight-tickets"
                      label="Weight tickets"
                      name="weightType"
                      value="weightTicket"
                      checked
                    />
                    <Field
                      as={Radio}
                      id="constructed-weight"
                      label="Constructed weight"
                      name="weightType"
                      value="constructedWeight"
                    />
                  </Fieldset>
                </FormGroup>
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
                <Label className={styles.label}>Net weight</Label>
                <div className={styles.displayValue}>4,565 lbs</div>
                <FormGroup>
                  <Fieldset>
                    <legend className="usa-label">Did they use a trailer they owned</legend>
                    <Field as={Radio} id="yes" label="Yes" name="ownTrailer" value="true" checked />
                    <Field as={Radio} id="no" label="No" name="ownTrailer" value="false" />
                  </Fieldset>
                </FormGroup>
                <h3 className={styles.reviewHeader}>Review trip {tripNumber}</h3>
                <p>Add a review for this weight ticket</p>
                <ApproveReject
                  id={mtoShipment?.ppmShipment?.id}
                  currentStatus={values.status}
                  rejectionReason={values.rejectionReason}
                  requestComplete
                  approvedStatus="APPROVED"
                  deniedStatus="DENIED"
                  canEditRejection={canEditRejection}
                  setCanEditRejection={setCanEditRejection}
                  handleApprovalChange={handleApprovalChange}
                  handleRejectChange={handleRejectChange}
                  handleRejectCancel={handleRejectCancel}
                  handleChange={handleChange}
                  handleFormReset={handleFormReset}
                />
              </Form>
            </div>
          );
        }}
      </Formik>
    </div>
  );
}

ReviewWeightTicket.propTypes = {
  mtoShipment: object.isRequired,
  tripNumber: string.isRequired,
  ppmNumber: string.isRequired,
};
