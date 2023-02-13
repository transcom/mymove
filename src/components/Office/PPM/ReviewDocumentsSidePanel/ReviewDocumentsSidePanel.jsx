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
  expenseTickets,
  proGearTickets,
  weightTickets,
}) {
  let status;
  let showReason;

  const handleSubmit = () => {
    // TODO: use a mutation query and then attach onSuccess and an onError handler
    onSuccess();
  };

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
    <Formik initialValues innerRef={formRef} onSubmit={handleSubmit}>
      <div className={classnames(styles.container, 'container--accent--ppm')}>
        <Form className={classnames(formStyles.form, styles.ReviewDocumentsSidePanel)}>
          <PPMHeaderSummary ppmShipment={ppmShipment} ppmNumber={ppmNumber} />
          <hr />
          <h3 className={styles.send}>Send to customer?</h3>
          <DocumentViewerSidebar.Content className={styles.sideBar}>
            {weightTickets.length > 0 && (
              <ul>
                {weightTickets.map((weight, index) => {
                  return (
                    <li className={styles.rowContainer} key={index}>
                      <div className={styles.row}>
                        <h3 className={styles.tripNumber}>Trip {index + 1}</h3>
                        {statusWithIcon(weight)}
                      </div>
                      {showReason ? <p>{weight.reason}</p> : null}
                    </li>
                  );
                })}
              </ul>
            )}
            {proGearTickets.length > 0 && (
              <ul>
                {proGearTickets.map((gear, index) => {
                  return (
                    <li className={styles.rowContainer} key={index}>
                      <div className={styles.row}>
                        <h3 className={styles.tripNumber}>Pro-gear {index + 1}</h3>
                        {statusWithIcon(gear)}
                      </div>
                      {showReason ? <p>{gear.reason}</p> : null}
                    </li>
                  );
                })}
              </ul>
            )}
            {expenseTickets.length > 0 && (
              <ul>
                {expenseTickets.map((exp, index) => {
                  return (
                    <li className={styles.rowContainer} key={index}>
                      <div className={styles.row}>
                        <h3 className={styles.tripNumber}>
                          {exp.movingExpenseType === expenseTypes.STORAGE ? 'Storage' : 'Receipt'}
                          &nbsp;{index + 1}
                        </h3>
                        {statusWithIcon(exp)}
                      </div>
                      {showReason ? <p>{exp.reason}</p> : null}
                    </li>
                  );
                })}
              </ul>
            )}
          </DocumentViewerSidebar.Content>
        </Form>
      </div>
    </Formik>
  );
}

ReviewDocumentsSidePanel.propTypes = {
  ppmShipment: PPMShipmentShape,
  ppmNumber: number,
  formRef: object,
  onSuccess: func,
  expenseTickets: array,
  proGearTickets: array,
  weightTickets: array,
};

ReviewDocumentsSidePanel.defaultProps = {
  ppmShipment: undefined,
  ppmNumber: 1,
  formRef: null,
  onSuccess: () => {},
  expenseTickets: [],
  proGearTickets: [],
  weightTickets: [],
};
