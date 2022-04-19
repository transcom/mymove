import { get, isEmpty } from 'lodash';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { push } from 'connected-react-router';
import { bindActionCreators } from 'redux';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import {
  approveReimbursement,
  downloadPPMAttachments,
  downloadPPMAttachmentsLabel,
  selectActivePPMForMove,
} from 'shared/Entities/modules/ppms';
import { selectAllDocumentsForMove } from 'shared/Entities/modules/moveDocuments';
import {
  getSignedCertification,
  selectPaymentRequestCertificationForMove,
} from 'shared/Entities/modules/signed_certifications';
import { getLastError } from 'shared/Swagger/selectors';
import { selectReimbursementById } from 'store/entities/selectors';

import { no_op } from 'shared/utils';
import { SIGNED_CERT_OPTIONS } from 'shared/constants';
import { formatCents } from 'shared/formatters';
import { formatDate } from 'utils/formatters';

import './PaymentsPanel.css';
import Alert from 'shared/Alert';
import ToolTip from 'shared/ToolTip';

const attachmentsErrorMessages = {
  422: 'Encountered an error while trying to create attachments bundle: Document is in the wrong format',
  424: 'Could not find any receipts or documents for this PPM',
  500: 'An unexpected error has occurred',
};

export function sswIsDisabled(ppm, signedCertification) {
  return missingSignature(signedCertification) || missingRequiredPPMInfo(ppm);
}

function missingSignature(signedCertification) {
  return isEmpty(signedCertification) || signedCertification.certification_type !== SIGNED_CERT_OPTIONS.PPM_PAYMENT;
}

function missingRequiredPPMInfo(ppm) {
  return isEmpty(ppm) || !ppm.actual_move_date || !ppm.net_weight;
}

function getUserDate() {
  return new Date().toISOString().split('T')[0];
}

// Taken from https://mathiasbynens.github.io/rel-noopener/
// tl;dr-- opening content in target _blank can leave parent window open to malicious code
// below is a safer way to open content in a new tab
function safeOpenInNewTab(url) {
  if (url) {
    let win = window.open();
    // win can be null if a pop-up blocker is used
    if (win) {
      win.opener = null;
      win.location = url;
    }
  }
}

class PaymentsTable extends Component {
  state = {
    showPaperwork: false,
    disableDownload: false,
  };

  componentDidMount() {
    const { moveId } = this.props;
    if (moveId != null) {
      this.props.getSignedCertification(moveId);
    }
  }

  approveReimbursement = () => {
    this.props.approveReimbursement(this.props.advance.id);
  };

  togglePaperwork = () => {
    this.setState((prevState) => ({ showPaperwork: !prevState.showPaperwork }));
  };

  disableDownloadAll = () => {
    return this.props.moveDocuments.length < 1;
  };

  startDownload = (docTypes) => {
    this.setState({ disableDownload: true });
    this.props.downloadPPMAttachments(this.props.ppm.id, docTypes).then((response) => {
      const {
        response: {
          obj: { url },
        },
      } = response;
      safeOpenInNewTab(url);
      this.setState({ disableDownload: false });
    });
  };

  documentUpload = () => {
    const { moveId } = this.props;
    this.props.push(`/moves/${moveId}/documents/new?move_document_type=SHIPMENT_SUMMARY`);
  };

  downloadShipmentSummary = () => {
    const { moveId } = this.props;
    const userDate = getUserDate();

    safeOpenInNewTab(`/internal/moves/${moveId}/shipment_summary_worksheet/?preparationDate=${userDate}`);
  };

  renderAdvanceAction = () => {
    if (this.props.ppm.status === 'APPROVED') {
      if (this.props.advance.status === 'APPROVED') {
        return <div>{/* Further actions to come*/}</div>;
      } else {
        return (
          <div onClick={this.approveReimbursement}>
            <ToolTip disabled={false} text="Approve" textStyle="tooltiptext-small">
              <FontAwesomeIcon aria-hidden className="icon approval-ready" icon="check" title="Approve" />
            </ToolTip>
          </div>
        );
      }
    } else {
      return (
        <ToolTip
          disabled={false}
          text={"Can't approve payment until shipment is approved"}
          textStyle="tooltiptext-medium"
        >
          <FontAwesomeIcon
            aria-hidden
            className="icon approval-blocked"
            icon="check"
            title="Can't approve payment until shipment is approved."
          />
        </ToolTip>
      );
    }
  };

  render() {
    const attachmentsError = this.props.attachmentsError;
    const advance = this.props.advance;
    const paperworkIcon = this.state.showPaperwork ? 'minus-square' : 'plus-square';

    return (
      <div className="payment-panel">
        <div className="payment-panel-title">Payments</div>
        <table className="payment-table">
          <tbody>
            <tr>
              <th className="payment-table-column-title" />
              <th className="payment-table-column-title">Amount</th>
              <th className="payment-table-column-title">Disbursement</th>
              <th className="payment-table-column-title">Requested on</th>
              <th className="payment-table-column-title">Status</th>
              <th className="payment-table-column-title">Actions</th>
            </tr>
            {!isEmpty(advance) ? (
              <React.Fragment>
                <tr>
                  <th className="payment-table-subheader" colSpan="6">
                    Payments against PPM Incentive
                  </th>
                </tr>
                <tr>
                  <td className="payment-table-column-content">Advance</td>
                  <td className="payment-table-column-content">
                    ${formatCents(get(advance, 'requested_amount')).toLocaleString()}
                  </td>
                  <td className="payment-table-column-content">{advance.method_of_receipt}</td>
                  <td className="payment-table-column-content">{formatDate(advance.requested_date)}</td>
                  <td className="payment-table-column-content">
                    {advance.status === 'APPROVED' ? (
                      <div>
                        <FontAwesomeIcon aria-hidden className="icon approval-ready" icon="check" title="Approved" />{' '}
                        Approved
                      </div>
                    ) : (
                      <div>
                        <FontAwesomeIcon
                          aria-hidden
                          className="icon approval-waiting"
                          icon="clock"
                          title="Awaiting Review"
                        />{' '}
                        Awaiting review
                      </div>
                    )}
                  </td>
                  <td className="payment-table-column-content">{this.renderAdvanceAction()}</td>
                </tr>
              </React.Fragment>
            ) : (
              <tr>
                <th className="payment-table-subheader">No payments requested</th>
              </tr>
            )}
          </tbody>
        </table>

        <div className="paperwork">
          <a onClick={this.togglePaperwork} className="usa-link">
            <FontAwesomeIcon aria-hidden className="icon" icon={paperworkIcon} />
            Create payment paperwork
          </a>
          {this.state.showPaperwork && (
            <Fragment>
              {attachmentsError && (
                <Alert type="error" heading="An error occurred">
                  {attachmentsErrorMessages[attachmentsError.statusCode] ||
                    'Something went wrong contacting the server.'}
                </Alert>
              )}
              <p>Complete the following steps in order to generate and file paperwork for payment:</p>
              <div className="paperwork">
                <div className="paperwork-step">
                  <div>
                    <p>Complete the Shipment Summary Worksheet</p>
                    <p>Open the SSW in a new tab. Download and open the PDF, then fill in any required info.</p>
                  </div>
                  <button
                    className="usa-button"
                    disabled={this.props.disableSSW}
                    onClick={this.downloadShipmentSummary}
                  >
                    Open SSW in a New Tab
                  </button>
                </div>
                {this.props.disableSSW && (
                  <Alert type="warning" heading="To get this worksheet ready to download:">
                    <ul>
                      <li>
                        Enter <b>departure date</b> (PPM tab)
                      </li>
                      <li>
                        Process <i>all</i> weight tickets to calculate <b>net weight</b>
                      </li>
                      <li>
                        Ask service member to <b>request payment</b>
                      </li>
                    </ul>
                    <p>
                      After that, if the button is still inactive, contact support at{' '}
                      <a href="tel:(628) 225-1540" className="usa-link">
                        (628) 225-1540
                      </a>
                    </p>
                  </Alert>
                )}
                <hr />

                <div className="paperwork-step">
                  <div>
                    <p>Download All Attachments (PDF)</p>
                    <p>Download bundle of PPM receipts and attach it to the completed Shipment Summary Worksheet.</p>
                  </div>
                  <button
                    className="usa-button"
                    disabled={this.state.disableDownload || this.disableDownloadAll()}
                    onClick={() =>
                      this.startDownload(['OTHER', 'WEIGHT_TICKET', 'WEIGHT_TICKET_SET', 'STORAGE_EXPENSE', 'EXPENSE'])
                    }
                  >
                    Download All Attachments (PDF)
                  </button>
                </div>

                <hr />

                <div className="paperwork-step">
                  <div>
                    <p>Download Orders and Weight Tickets (PDF)</p>
                    <p>
                      Download bundle of Orders and Weight Tickets (without receipts) and attach it to the completed
                      Shipment Summary Worksheet.
                    </p>
                  </div>
                  <button
                    className="usa-button"
                    disabled={this.state.disableDownload}
                    onClick={() =>
                      this.startDownload(['OTHER', 'WEIGHT_TICKET', 'WEIGHT_TICKET_SET', 'STORAGE_EXPENSE'])
                    }
                  >
                    Download Orders and Weight Tickets (PDF)
                  </button>
                </div>

                <hr />

                <div className="paperwork-step">
                  <div>
                    <p>Upload completed packet</p>
                    <p>
                      Save the worksheet and attachments together as one PDF. Then upload the completed packet for
                      customer and Finance.
                    </p>
                  </div>
                  <button className="usa-button" onClick={this.documentUpload}>
                    Upload Completed Packet
                  </button>
                </div>
              </div>
            </Fragment>
          )}
        </div>
      </div>
    );
  }
}

const mapStateToProps = (state, ownProps) => {
  const { moveId } = ownProps;
  const ppm = selectActivePPMForMove(state, moveId);
  const advance = selectReimbursementById(state, ppm.advance) || {};
  const signedCertifications = selectPaymentRequestCertificationForMove(state, moveId);
  const moveDocuments = selectAllDocumentsForMove(state, moveId);
  const disableSSW = sswIsDisabled(ppm, signedCertifications);
  return {
    ppm,
    disableSSW,
    moveId,
    advance,
    attachmentsError: getLastError(state, `${downloadPPMAttachmentsLabel}-${moveId}`),
    moveDocuments,
  };
};

const mapDispatchToProps = (dispatch) =>
  bindActionCreators(
    {
      getSignedCertification,
      approveReimbursement,
      update: no_op,
      downloadPPMAttachments,
      push,
    },
    dispatch,
  );

export default connect(mapStateToProps, mapDispatchToProps)(PaymentsTable);
