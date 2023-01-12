import React from 'react';
import { Form } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { number } from 'prop-types';

import PPMHeaderSummary from '../PPMHeaderSummary/PPMHeaderSummary';

import styles from './ReviewDocumentsSidePanel.module.scss';

import { PPMShipmentShape } from 'types/shipment';
import formStyles from 'styles/form.module.scss';
import DocumentViewerSidebar from 'pages/Office/DocumentViewerSidebar/DocumentViewerSidebar';
import { ReviewDocumentsStatus } from 'constants/ppms';
import { expenseTypes } from 'constants/ppmExpenseTypes';

export default function ReviewDocumentsSidePanel({ ppmShipment, ppmNumber, expenseNumber, tripNumber }) {
  // TODO: return pro-gear tickets & expenses data will be done in later tickets. For now, we have placeholders here (proGearTickets, expenseTickets).

  const proGearTickets = [
    { status: ReviewDocumentsStatus.EXCLUDE, reason: 'Objects not applicable' },
    { status: ReviewDocumentsStatus.ACCEPT, reason: null },
  ];

  const expenseTickets = [
    { movingExpenseType: expenseTypes.STORAGE, status: ReviewDocumentsStatus.REJECT, reason: 'Too large' },
    { movingExpenseType: expenseTypes.PACKING_MATERIALS, status: ReviewDocumentsStatus.ACCEPT, reason: null },
  ];

  // TODO: ability to reject/approve weight tickets will be done in another ticker. For now, this is placeholder for an accepted weight ticket

  const weightTickets = [{ status: ReviewDocumentsStatus.ACCEPT, reason: null }];

  let status;
  let showReason;

  const statusWithIcon = (ticket) => {
    if (ticket.status === ReviewDocumentsStatus.ACCEPT) {
      status = (
        <div className={styles.iconRow}>
          <FontAwesomeIcon icon="check" />
          <span>Accept</span>
        </div>
      );
    } else if (ticket.status === ReviewDocumentsStatus.EXCLUDE) {
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
      <Form className={classnames(formStyles.form, styles.ReviewDocumentsSidePanel)}>
        <PPMHeaderSummary ppmShipment={ppmShipment} ppmNumber={ppmNumber} />
        <hr />
        <h3 className={styles.send}>Send to customer?</h3>
        <DocumentViewerSidebar.Content>
          {weightTickets.length > 0
            ? weightTickets.map((weight) => {
                return (
                  <div className={styles.rowContainer}>
                    <div className={styles.row}>
                      <h3 className={styles.tripNumber}>Trip {tripNumber}</h3>
                      {statusWithIcon(weight)}
                    </div>
                    {showReason ? <p>{weight.weightTicketReason}</p> : null}
                  </div>
                );
              })
            : null}
          {proGearTickets.length > 0
            ? proGearTickets.map((gear, index) => {
                return (
                  <div className={styles.rowContainer}>
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
            ? expenseTickets.map((exp) => {
                return (
                  <div className={styles.rowContainer}>
                    <div className={styles.row}>
                      <h3 className={styles.tripNumber}>
                        {exp.movingExpenseType === expenseTypes.STORAGE ? 'Storage' : 'Receipt'}
                        &nbsp;{expenseNumber}
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
    </div>
  );
}

ReviewDocumentsSidePanel.propTypes = {
  ppmShipment: PPMShipmentShape,
  tripNumber: number,
  ppmNumber: number,
  expenseNumber: number,
};

ReviewDocumentsSidePanel.defaultProps = {
  ppmShipment: undefined,
  tripNumber: 1,
  ppmNumber: 1,
  expenseNumber: 1,
};
