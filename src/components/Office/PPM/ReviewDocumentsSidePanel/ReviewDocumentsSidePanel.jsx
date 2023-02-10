import React from 'react';
import { Form } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { array, func, number, object } from 'prop-types';

import PPMHeaderSummary from '../PPMHeaderSummary/PPMHeaderSummary';

import styles from './ReviewDocumentsSidePanel.module.scss';

import { PPMShipmentShape } from 'types/shipment';
import formStyles from 'styles/form.module.scss';
import DocumentViewerSidebar from 'pages/Office/DocumentViewerSidebar/DocumentViewerSidebar';
import PPMDocumentsStatus from 'constants/ppms';
import { expenseTypes } from 'constants/ppmExpenseTypes';

export default function ReviewDocumentsSidePanel({
  ppmShipment,
  ppmNumber,
  formRef,
  onSuccess,
  onError,
  expenseTickets,
  proGearTickets,
  weightTickets,
}) {
  // TODO: return pro-gear tickets & expenses data will be done in later tickets. For now, we have placeholders here (proGearTickets, expenseTickets).

  // TODO: ability to reject/approve weight tickets will be done in another ticker. For now, this is placeholder for an accepted weight ticket

  // TODO: will need something like this for form submission

  // const [patchWeightTicketMutation] = useMutation(patchWeightTicket, {
  //   onSuccess,
  //   onError,
  // });

  // const [confirmPPMDocuments] = useMutation(confirmPPMDocuments, {
  //   onSuccess,
  //   onError,
  // });

  const handleSubmit = () => {
    // TODO: use the mutation and pass in onSuccess and onError, like ReviewWeightTicket
    if (Math.random() < 0.9999) {
      onSuccess();
    } else {
      onError();
    }
  };

  let status;
  let showReason;

  const statusWithIcon = (ticket) => {
    if (ticket.status === PPMDocumentsStatus.APPROVED) {
      status = (
        <div className={styles.iconRow}>
          <FontAwesomeIcon icon="check" />
          <span>Accept</span>
        </div>
      );
    } else if (ticket.status === PPMDocumentsStatus.EXCLUDED) {
      status = (
        <div className={styles.iconRow}>
          <FontAwesomeIcon icon="ban" />
          <span>Exclude</span>
        </div>
      );
      showReason = true;
    } else {
      status = (
        <div className={styles.iconRow}>
          <FontAwesomeIcon icon="times" />
          <span>Reject</span>
        </div>
      );
      showReason = true;
    }
    return status;
  };

  return (
    <div className={classnames(styles.container, 'container--accent--ppm')}>
      <Formik innerRef={formRef} onSubmit={handleSubmit} initialValues>
        <Form className={classnames(formStyles.form, styles.ReviewDocumentsSidePanel)}>
          <PPMHeaderSummary ppmShipment={ppmShipment} ppmNumber={ppmNumber} />
          <hr />
          <h3 className={styles.send}>Send to customer?</h3>
          <DocumentViewerSidebar.Content className={styles.sideBar}>
            {weightTickets?.length > 0
              ? weightTickets.map((weight, index) => {
                  return (
                    <div className={styles.rowContainer} key={index}>
                      <div className={styles.row}>
                        <h3 className={styles.tripNumber}>Trip {index + 1}</h3>
                        {statusWithIcon(weight)}
                      </div>
                      {showReason ? <p>{weight.reason}</p> : null}
                    </div>
                  );
                })
              : null}
            {proGearTickets.length > 0
              ? proGearTickets.map((gear, index) => {
                  return (
                    <div className={styles.rowContainer} key={index}>
                      <div className={styles.row}>
                        <h3 className={styles.tripNumber}>Pro-gear {index + 1}</h3>
                        {statusWithIcon(gear)}
                      </div>
                      {showReason ? <p>{gear.reason}</p> : null}
                    </div>
                  );
                })
              : null}
            {expenseTickets.length > 0
              ? expenseTickets.map((exp, index) => {
                  return (
                    <div className={styles.rowContainer} key={index}>
                      <div className={styles.row}>
                        <h3 className={styles.tripNumber}>
                          {exp.movingExpenseType === expenseTypes.STORAGE ? 'Storage' : 'Receipt'}
                          &nbsp;{index + 1}
                        </h3>
                        {statusWithIcon(exp)}
                      </div>
                      {showReason ? <p>{exp.reason}</p> : null}
                    </div>
                  );
                })
              : null}
          </DocumentViewerSidebar.Content>
        </Form>
      </Formik>
    </div>
  );
}

ReviewDocumentsSidePanel.propTypes = {
  ppmShipment: PPMShipmentShape,
  ppmNumber: number,
  formRef: object,
  onSuccess: func,
  onError: func,
  expenseTickets: array,
  proGearTickets: array,
  weightTickets: array,
};

ReviewDocumentsSidePanel.defaultProps = {
  ppmShipment: undefined,
  ppmNumber: 1,
  formRef: null,
  onSuccess: () => {},
  onError: () => {},
  expenseTickets: [],
  proGearTickets: [],
  weightTickets: [],
};
