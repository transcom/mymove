import React, { Component } from 'react';
import { bool } from 'prop-types';
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import { selectPPMCloseoutDocumentsForMove } from 'shared/Entities/modules/movingExpenseDocuments';
import { deleteMoveDocument, getMoveDocumentsForMove } from 'shared/Entities/modules/moveDocuments';
import docsAddedCheckmarkImg from 'shared/images/docs_added_checkmark.png';
import WeightTicketListItem from './WeightTicketListItem';
import ExpenseTicketListItem from './ExpenseTicketListItem';
import { formatExpenseDocs } from '../utility';

import './PaymentReview.css';

export class DocumentsUploaded extends Component {
  state = {
    showDocs: false,
  };

  static propTypes = {
    showLinks: bool,
    inReviewPage: bool,
  };

  static defaultProps = {
    showLinks: false,
    inReviewPage: false,
  };

  componentDidMount() {
    const { moveId } = this.props;
    this.props.getMoveDocumentsForMove(moveId);
  }

  toggleShowDocs = () => {
    this.setState((prevState) => ({ showDocs: !prevState.showDocs }));
  };

  renderHeader = () => {
    const { expenseDocs, weightTicketSetDocs, weightTicketDocs, inReviewPage } = this.props;
    const totalDocs = expenseDocs.length + weightTicketSetDocs.length + weightTicketDocs.length;
    const documentLabel = `document${totalDocs > 1 ? 's' : ''}`;

    return <h3>{inReviewPage ? `Document Summary - ${totalDocs} total` : `${totalDocs} ${documentLabel} added`}</h3>;
  };

  render() {
    const { showDocs } = this.state;
    const { expenseDocs, weightTicketSetDocs, weightTicketDocs, moveId, showLinks, inReviewPage, deleteMoveDocument } =
      this.props;
    const totalDocs = expenseDocs.length + weightTicketSetDocs.length + weightTicketDocs.length;
    const expandedDocumentList = showDocs || inReviewPage;
    const hiddenDocumentList = !inReviewPage && !showDocs;

    if (totalDocs === 0) {
      return null;
    }
    return (
      <div
        className="doc-summary-container"
        data-testid="documents-uploaded"
        style={{ paddingBottom: hiddenDocumentList ? '1em' : null, marginTop: !inReviewPage ? '1em' : null }}
      >
        <div className="documents-uploaded-header">
          {!inReviewPage && (
            <img
              alt="documents added checkmark"
              src={docsAddedCheckmarkImg}
              style={{ alignSelf: 'center', marginRight: 5 }}
            />
          )}
          {this.renderHeader()}
          {!inReviewPage && (
            <a
              data-testid="toggle-documents-uploaded"
              style={{ paddingLeft: '1em' }}
              onClick={this.toggleShowDocs}
              className="usa-link"
            >
              {showDocs ? 'Hide' : 'Show'}
            </a>
          )}
        </div>
        {expandedDocumentList && (
          <>
            {weightTicketDocs.length > 0 && (
              <>
                <h4>{weightTicketDocs.length} weight tickets</h4>
                <div className="tickets">
                  {weightTicketDocs.map((ticket, index) => (
                    <WeightTicketListItem
                      key={ticket.id}
                      num={index}
                      showDelete={inReviewPage}
                      deleteDocumentListItem={deleteMoveDocument}
                      isWeightTicketSet={false}
                      {...ticket}
                    />
                  ))}
                </div>
                <hr id="doc-summary-separator" />
              </>
            )}
            <h4>{weightTicketSetDocs.length} sets of weight tickets</h4>
            <div className="tickets">
              {weightTicketSetDocs.map((ticket, index) => (
                <WeightTicketListItem
                  key={ticket.id}
                  num={index}
                  showDelete={inReviewPage}
                  deleteDocumentListItem={deleteMoveDocument}
                  isWeightTicketSet={true}
                  uploads={ticket.document.uploads}
                  {...ticket}
                />
              ))}
            </div>
            {showLinks && (
              <Link data-testid="weight-ticket-link" to={`/moves/${moveId}/ppm-weight-ticket`} className="usa-link">
                <FontAwesomeIcon className="icon link-blue" icon="plus-circle" /> Add weight ticket
              </Link>
            )}
            <hr id="doc-summary-separator" />
            <h4 style={{ paddingBottom: expenseDocs.length === 0 ? '1em' : null }}>
              {expenseDocs.length} expense{expenseDocs.length >= 0 ? 's' : ''}
            </h4>
            <div className="tickets">
              {formatExpenseDocs(expenseDocs).map((expense) => (
                <ExpenseTicketListItem
                  key={expense.id}
                  showDelete={inReviewPage}
                  deleteDocumentListItem={deleteMoveDocument}
                  uploads={expense.uploads}
                  {...expense}
                />
              ))}
            </div>
            {showLinks && (
              <div className="add-expense-link">
                <Link data-testid="expense-link" to={`/moves/${moveId}/ppm-expenses`} className="usa-link">
                  <FontAwesomeIcon className="icon link-blue" icon="plus-circle" /> Add expense
                </Link>
              </div>
            )}
          </>
        )}
      </div>
    );
  }
}

function mapStateToProps(state, { moveId }) {
  return {
    moveId,
    expenseDocs: selectPPMCloseoutDocumentsForMove(state, moveId, ['EXPENSE']),
    weightTicketSetDocs: selectPPMCloseoutDocumentsForMove(state, moveId, ['WEIGHT_TICKET_SET']),
    weightTicketDocs: selectPPMCloseoutDocumentsForMove(state, moveId, ['WEIGHT_TICKET']),
  };
}

const mapDispatchToProps = {
  selectPPMCloseoutDocumentsForMove,
  getMoveDocumentsForMove,
  deleteMoveDocument,
};

export default connect(mapStateToProps, mapDispatchToProps)(DocumentsUploaded);
