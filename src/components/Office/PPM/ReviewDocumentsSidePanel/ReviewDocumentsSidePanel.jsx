import React from 'react';
import { useMutation } from '@tanstack/react-query';
import { Form } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { arrayOf, func, number, object } from 'prop-types';
import moment from 'moment';

import PPMHeaderSummary from '../PPMHeaderSummary/PPMHeaderSummary';

import styles from './ReviewDocumentsSidePanel.module.scss';

import { expenseTypes } from 'constants/ppmExpenseTypes';
import { OrderShape } from 'types/order';
import { patchPPMDocumentsSetStatus } from 'services/ghcApi';
import { ExpenseShape, PPMShipmentShape, ProGearTicketShape, WeightTicketShape } from 'types/shipment';
import formStyles from 'styles/form.module.scss';
import DocumentViewerSidebar from 'pages/Office/DocumentViewerSidebar/DocumentViewerSidebar';
import PPMDocumentsStatus from 'constants/ppms';
import { formatDate, formatCents } from 'utils/formatters';

export default function ReviewDocumentsSidePanel({
  ppmShipment,
  ppmShipmentInfo,
  ppmNumber,
  formRef,
  onSuccess,
  onError,
  expenseTickets,
  proGearTickets,
  weightTickets,
  readOnly,
  order,
}) {
  let status;
  let showReason;

  const { mutate: patchDocumentsSetStatusMutation } = useMutation(patchPPMDocumentsSetStatus, {
    onSuccess,
    onError,
  });

  const handleSubmit = () => {
    if (readOnly) {
      onSuccess();
      return;
    }
    patchDocumentsSetStatusMutation({
      ppmShipmentId: ppmShipment.id,
      eTag: ppmShipment.eTag,
    });
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

  const expenseSetProjection = (expenses) => {
    const process = expenses.reduce((accumulator, item, index) => {
      accumulator[item.movingExpenseType] ??= [];
      const expenseSet = accumulator[item.movingExpenseType];

      expenseSet.push({ ...item, receiptIndex: index + 1, groupIndex: expenseSet.length + 1 });
      return accumulator;
    }, {});

    const compareReceiptIndex = (itemA, itemB) => itemA.receiptIndex >= itemB.receiptIndex;

    return Object.values(process).flat().sort(compareReceiptIndex);
  };

  const titleCase = (input) => input.charAt(0).toUpperCase() + input.slice(1);
  const allCase = (input) => input.split(' ').map(titleCase).join(' ');
  const formatMovingType = (input) => allCase(input?.trim().toLowerCase().replace('_', ' ') ?? '');
  let total = 0;

  return (
    <Formik initialValues innerRef={formRef} onSubmit={handleSubmit}>
      <div className={classnames(styles.container, 'container--accent--ppm')}>
        <div className={classnames(formStyles.form, styles.ReviewDocumentsSidePanel, styles.PPMHeaderSummary)}>
          <PPMHeaderSummary
            ppmShipmentInfo={ppmShipmentInfo}
            order={order}
            ppmNumber={ppmNumber}
            showAllFields
            readOnly={readOnly}
          />
        </div>
        <Form className={classnames(formStyles.form, styles.ReviewDocumentsSidePanel)}>
          <hr />
          <h3 className={styles.send}>{readOnly ? 'Sent to customer' : 'Send to customer?'}</h3>
          <DocumentViewerSidebar.Content className={styles.sideBar}>
            <ul>
              {weightTickets.length > 0
                ? weightTickets.map((weight, index) => {
                    return (
                      <li className={styles.rowContainer} key={index}>
                        <div className={styles.row}>
                          <h3 className={styles.tripNumber}>Trip {index + 1}</h3>
                          {statusWithIcon(weight)}
                        </div>
                        {showReason ? <p>{weight.reason}</p> : null}

                        <dl className={classnames(styles.ItemDetails)}>
                          <span>
                            <dt>Empty Weight:</dt>
                            <dd>{weight.emptyWeight} lbs</dd>
                          </span>
                          <span>
                            <dt>Full Weight:</dt>
                            <dl>{weight.fullWeight} lbs</dl>
                          </span>
                          <span>
                            <dt>Net Weight:</dt>
                            <dl>{weight.fullWeight - weight.emptyWeight} lbs</dl>
                          </span>
                          <span>
                            <dt>Trailer Used:</dt>
                            <dl>{weight.ownsTrailer ? `Yes` : `No`}</dl>
                          </span>
                          {weight.ownsTrailer && (
                            <span>
                              <dt>Trailer Claimable:</dt>
                              <dl>{weight.trailerMeetsCriteria ? `Yes` : `No`}</dl>
                            </span>
                          )}
                        </dl>
                      </li>
                    );
                  })
                : null}
              {proGearTickets.length > 0
                ? proGearTickets.map((gear, index) => {
                    return (
                      <li className={styles.rowContainer} key={index}>
                        <div className={styles.row}>
                          <h3 className={styles.tripNumber}>Pro-gear {index + 1}</h3>
                          {statusWithIcon(gear)}
                        </div>
                        {showReason ? <p>{gear.reason}</p> : null}

                        <dl className={classnames(styles.ItemDetails)}>
                          <span>
                            <dt>Belongs To: </dt>
                            <dd>{gear.belongsToSelf ? `Customer` : `Spouse`}</dd>
                          </span>
                          <span>
                            <dt>Missing Weight Ticket (Constructed)?</dt>
                            <dl>{gear.missingWeightTicket ? `Yes` : `No`}</dl>
                          </span>
                          <span>
                            <dt>Pro-gear Weight:</dt>
                            {/* TODO: proGearWeight shows empty for some reason? */}
                            <dl>{gear.weight} lbs</dl>
                          </span>
                        </dl>
                      </li>
                    );
                  })
                : null}
              {expenseTickets.length > 0
                ? expenseSetProjection(expenseTickets).map((exp) => {
                    if (exp.status === PPMDocumentsStatus.APPROVED) {
                      if (exp.movingExpenseType === expenseTypes.STORAGE) {
                        total += exp.sitReimburseableAmount;
                      } else {
                        total += exp.amount;
                      }
                    }
                    return (
                      <li className={styles.rowContainer} key={exp.receiptIndex}>
                        <div className={styles.row}>
                          <h3 className={styles.tripNumber}>
                            Receipt&nbsp;{exp.receiptIndex}
                            <br />
                            {formatMovingType(exp.movingExpenseType)}&nbsp;#{exp.groupIndex}
                          </h3>
                          {statusWithIcon(exp)}
                        </div>
                        {showReason ? <p>{exp.reason}</p> : null}

                        <div className={classnames(styles.ItemDetails)}>
                          {exp.movingExpenseType === expenseTypes.STORAGE ? (
                            <dl>
                              <span>
                                <dt>SIT Start Date:</dt>
                                <dd>{formatDate(exp.sitStartDate)}</dd>
                              </span>
                              <span>
                                <dt>SIT End Date:</dt>
                                <dl>{formatDate(exp.sitEndDate)}</dl>
                              </span>
                              <span>
                                <dt>Total Days in SIT:</dt>
                                <dl data-testid="days-in-sit">
                                  {moment(exp.sitEndDate, 'YYYY MM DD')
                                    .add(1, 'days')
                                    .diff(moment(exp.sitStartDate, 'YYYY MM DD'), 'days')}
                                </dl>
                              </span>
                              <span>
                                <dt>Authorized Price:</dt>
                                <dl>${formatCents(exp.sitReimburseableAmount)}</dl>
                              </span>
                            </dl>
                          ) : (
                            <span>
                              <dt>Authorized Price:</dt>
                              <dl>${formatCents(exp.amount)}</dl>
                            </span>
                          )}
                        </div>
                      </li>
                    );
                  })
                : null}
              {expenseTickets.length > 0 ? (
                <>
                  <hr />
                  <li className={styles.rowContainer}>
                    <div className={classnames(styles.ItemDetails)}>
                      <dl>
                        <span className={classnames(styles.ReceiptTotal)}>
                          <dt>Accepted Receipt Totals:</dt>
                          <dd>${formatCents(total)}</dd>
                        </span>
                      </dl>
                    </div>
                  </li>
                </>
              ) : null}
            </ul>
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
  onError: func,
  expenseTickets: arrayOf(ExpenseShape),
  proGearTickets: arrayOf(ProGearTicketShape),
  weightTickets: arrayOf(WeightTicketShape),
  order: OrderShape.isRequired,
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
