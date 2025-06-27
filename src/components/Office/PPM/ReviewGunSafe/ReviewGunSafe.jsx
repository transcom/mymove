import React, { useState, useEffect } from 'react';
import { useMutation } from '@tanstack/react-query';
import { func, number, string, object } from 'prop-types';
import { Field, Formik } from 'formik';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Form, FormGroup, Label, Radio, Button, Textarea } from '@trussworks/react-uswds';
import * as Yup from 'yup';

import PPMHeaderSummary from '../PPMHeaderSummary/PPMHeaderSummary';

import styles from './ReviewGunSafe.module.scss';

import { ErrorMessage } from 'components/form';
import { OrderShape } from 'types/order';
import { patchGunSafeWeightTicket } from 'services/ghcApi';
import { GunSafeTicketShape } from 'types/shipment';
import Fieldset from 'shared/Fieldset';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import formStyles from 'styles/form.module.scss';
import approveRejectStyles from 'styles/approveRejectControls.module.scss';
import ppmDocumentStatus from 'constants/ppms';
import Hint from 'components/Hint';

const validationSchema = Yup.object().shape({
  gunSafeWeight: Yup.number()
    .min(0, 'Enter a weight 0 lbs or greater')
    .max(500, 'Enter a weight 500 lbs or less')
    .when('missingWeightTicket', {
      is: 'true',
      then: (schema) => schema.required('Enter the constructed gun safe weight'),
      otherwise: (schema) => schema.required('Enter the weight with gun safe'),
    }),
  description: Yup.string().required('Required'),
  missingWeightTicket: Yup.string(),
  status: Yup.string().required('Reviewing this gun safe is required'),
  rejectionReason: Yup.string().when('status', {
    is: ppmDocumentStatus.REJECTED,
    then: (schema) => schema.required('Add a reason why this gun safe is rejected'),
  }),
});

export default function ReviewGunSafe({
  ppmShipmentInfo,
  gunSafe,
  tripNumber,
  ppmNumber,
  onError,
  onSuccess,
  formRef,
  readOnly,
  order,
}) {
  const [canEditRejection, setCanEditRejection] = useState(true);

  const { mutate: patchGunSafeMutation } = useMutation(patchGunSafeWeightTicket, {
    onSuccess,
    onError,
  });

  const { description, hasWeightTickets, weight, status, reason } = gunSafe || {};

  const missingWeightTicketValue = hasWeightTickets ? 'false' : 'true';

  const handleSubmit = (values) => {
    if (readOnly) {
      onSuccess();
      return;
    }

    let hasWeightTicketValue;
    if (values.missingWeightTicket === 'true') {
      hasWeightTicketValue = false;
    }
    if (values.missingWeightTicket === 'false') {
      hasWeightTicketValue = true;
    }
    const payload = {
      ppmShipmentId: gunSafe.ppmShipmentId,
      hasWeightTickets: hasWeightTicketValue,
      weight: parseInt(values.gunSafeWeight, 10),
      reason: values.status === ppmDocumentStatus.APPROVED ? null : values.rejectionReason,
      status: values.status,
    };
    patchGunSafeMutation({
      ppmShipmentId: gunSafe.ppmShipmentId,
      gunSafeWeightTicketId: gunSafe.id,
      payload,
      eTag: gunSafe.eTag,
    });
  };

  const initialValues = {
    status: status || '',
    rejectionReason: reason || '',
    missingWeightTicket: missingWeightTicketValue,
    description: description ? `${description}` : '',
    gunSafeWeight: weight ? `${weight}` : '',
  };

  useEffect(() => {
    if (formRef?.current) {
      formRef.current.resetForm();
      formRef.current.validateForm();
    }
  }, [formRef, gunSafe]);

  return (
    <div className={classnames(styles.container, 'container--accent--ppm')}>
      <Formik
        initialValues={initialValues}
        validationSchema={validationSchema}
        innerRef={formRef}
        onSubmit={handleSubmit}
        enableReinitialize
        validateOnMount
      >
        {({ handleChange, errors, touched, values }) => {
          const handleApprovalChange = (event) => {
            handleChange(event);
            setCanEditRejection(true);
          };

          return (
            <>
              <div className={classnames(formStyles.form, styles.reviewGunSafe, styles.headerContainer)}>
                <PPMHeaderSummary
                  ppmShipmentInfo={ppmShipmentInfo}
                  order={order}
                  ppmNumber={ppmNumber}
                  showAllFields={false}
                  readOnly={readOnly}
                />
              </div>
              <Form className={classnames(formStyles.form, styles.reviewGunSafe)}>
                <hr />
                <h3 className={styles.tripNumber}>Gun safe {tripNumber}</h3>
                <legend className={classnames('usa-label', styles.label)}>Description</legend>
                <div className={styles.displayValue}>{values.description}</div>
                <FormGroup>
                  <Fieldset>
                    <legend className="usa-label">Gun safe weight</legend>
                    <Field
                      as={Radio}
                      id="weight-tickets"
                      label="Weight tickets"
                      data-testid="gunSafeWeightTicket"
                      name="missingWeightTicket"
                      value="false"
                      checked={values.missingWeightTicket === 'false'}
                      disabled={readOnly}
                    />
                    <Field
                      as={Radio}
                      id="constructed-weight"
                      label="Constructed weight"
                      data-testid="gunSafeConstructedWeight"
                      name="missingWeightTicket"
                      value="true"
                      checked={values.missingWeightTicket === 'true'}
                      disabled={readOnly}
                    />
                  </Fieldset>
                </FormGroup>
                <MaskedTextField
                  defaultValue="0"
                  name="gunSafeWeight"
                  label={
                    values.missingWeightTicket === 'true' ? 'Constructed gun safe weight' : "Shipment's gun safe weight"
                  }
                  id="gunSafeWeight"
                  data-testid="gunSafeWeight"
                  mask={Number}
                  scale={0} // digits after point, 0 for integers
                  thousandsSeparator=","
                  lazy={false} // immediate masking evaluation
                  suffix="lbs"
                  disabled={readOnly}
                />
                {values.gunSafeWeight > 500 && (
                  <Hint>
                    The government authorizes the shipment of a gun safe up to 500 lbs (This is not charged against the
                    authorized weight entitlement. Any weight over 500 lbs is charged against the weight entitlement).
                  </Hint>
                )}
                <h3 className={styles.reviewHeader}>Review gun safe {tripNumber}</h3>
                <p>Add a review for this gun safe</p>
                <ErrorMessage display={!!errors?.status && !!touched?.status}>{errors.status}</ErrorMessage>
                <Fieldset>
                  <div
                    className={classnames(approveRejectStyles.statusOption, {
                      [approveRejectStyles.selected]: values.status === ppmDocumentStatus.APPROVED,
                    })}
                  >
                    <Radio
                      id={`approve-${gunSafe?.id}`}
                      checked={values.status === ppmDocumentStatus.APPROVED}
                      value={ppmDocumentStatus.APPROVED}
                      name="status"
                      label="Accept"
                      onChange={handleApprovalChange}
                      data-testid="approveRadio"
                      className={styles.acceptRadio}
                      disabled={readOnly}
                    />
                  </div>
                  <div
                    className={classnames(approveRejectStyles.statusOption, styles.reject, {
                      [approveRejectStyles.selected]: values.status === ppmDocumentStatus.REJECTED,
                    })}
                  >
                    <Radio
                      id={`reject-${gunSafe?.id}`}
                      checked={values.status === ppmDocumentStatus.REJECTED}
                      value={ppmDocumentStatus.REJECTED}
                      name="status"
                      label="Reject"
                      onChange={handleChange}
                      data-testid="rejectRadio"
                      className={styles.rejectRadio}
                      disabled={readOnly}
                    />

                    {values.status === ppmDocumentStatus.REJECTED && (
                      <FormGroup className={styles.rejectionReason}>
                        <Label htmlFor={`rejectReason-${gunSafe?.id}`}>Reason</Label>
                        {!canEditRejection && (
                          <>
                            <p data-testid="rejectionReasonReadOnly">{gunSafe?.reason || values.rejectionReason}</p>
                            <Button
                              type="button"
                              unstyled
                              data-testid="editReasonButton"
                              className={styles.clearStatus}
                              onClick={() => setCanEditRejection(true)}
                              aria-label="Edit reason button"
                              disabled={readOnly}
                            >
                              <span className="icon">
                                <FontAwesomeIcon icon="pen" title="Edit reason" alt="" />
                              </span>
                              <span aria-hidden="true">Edit reason</span>
                            </Button>
                          </>
                        )}

                        {canEditRejection && (
                          <>
                            <ErrorMessage display={!!errors?.rejectionReason && !!touched?.rejectionReason}>
                              {errors.rejectionReason}
                            </ErrorMessage>
                            <Textarea
                              id={`rejectReason-${gunSafe?.id}`}
                              name="rejectionReason"
                              onChange={handleChange}
                              error={touched.rejectionReason ? errors.rejectionReason : null}
                              value={values.rejectionReason}
                              placeholder="Type something"
                              disabled={readOnly}
                            />
                            <div className={styles.hint}>{500 - values.rejectionReason.length} characters</div>
                          </>
                        )}
                      </FormGroup>
                    )}
                  </div>
                </Fieldset>
              </Form>
            </>
          );
        }}
      </Formik>
    </div>
  );
}

ReviewGunSafe.propTypes = {
  gunSafe: GunSafeTicketShape,
  tripNumber: number.isRequired,
  ppmNumber: string.isRequired,
  onSuccess: func,
  formRef: object,
  order: OrderShape.isRequired,
};

ReviewGunSafe.defaultProps = {
  gunSafe: null,
  onSuccess: null,
  formRef: null,
};
