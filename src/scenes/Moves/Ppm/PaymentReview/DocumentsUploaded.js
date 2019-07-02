import React, { Component } from 'react';
import { bool } from 'prop-types';
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPlusCircle from '@fortawesome/fontawesome-free-solid/faPlusCircle';
import { selectPPMCloseoutDocumentsForMove } from 'shared/Entities/modules/movingExpenseDocuments';
import { getMoveDocumentsForMove } from 'shared/Entities/modules/moveDocuments';
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
    this.setState({ showDocs: !this.state.showDocs });
  };

  renderHeader = () => {
    const { expenseDocs, weightTicketDocs, inReviewPage } = this.props;
    const totalDocs = expenseDocs.length + weightTicketDocs.length;
    const documentLabel = `document${totalDocs > 1 ? 's' : ''}`;

    return <h3>{inReviewPage ? `Document Summary - ${totalDocs} total` : `${totalDocs} ${documentLabel} added`}</h3>;
  };

  render() {
    const { showDocs } = this.state;
    const { expenseDocs, weightTicketDocs, moveId, showLinks, inReviewPage } = this.props;
    const totalDocs = expenseDocs.length + weightTicketDocs.length;
    const expandedDocumentList = showDocs || inReviewPage;
    const hiddenDocumentList = !inReviewPage && !showDocs;

    if (totalDocs === 0) {
      return null;
    }
    return (
      <div
        className="doc-summary-container"
        data-cy="documents-uploaded"
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
            <a data-cy="toggle-documents-uploaded" style={{ paddingLeft: '1em' }} onClick={this.toggleShowDocs}>
              {showDocs ? 'Hide' : 'Show'}
            </a>
          )}
        </div>
        {expandedDocumentList && (
          <>
            <h4>{weightTicketDocs.length} sets of weight tickets</h4>
            <div className="tickets">
              {weightTicketDocs.map((ticket, index) => (
                <WeightTicketListItem key={ticket.id} num={index} showDelete={inReviewPage} {...ticket} />
              ))}
            </div>
            {showLinks && (
              <Link data-cy="weight-ticket-link" to={`/moves/${moveId}/ppm-weight-ticket`}>
                <FontAwesomeIcon className="icon link-blue" icon={faPlusCircle} /> Add weight ticket
              </Link>
            )}
            <hr id="doc-summary-separator" />
            <h4 style={{ paddingBottom: expenseDocs.length === 0 ? '1em' : null }}>
              {expenseDocs.length} expense{expenseDocs.length >= 0 ? 's' : ''}
            </h4>
            <div className="tickets">
              {formatExpenseDocs(expenseDocs).map(expense => (
                <ExpenseTicketListItem key={expense.id} showDelete={inReviewPage} {...expense} />
              ))}
            </div>
            {showLinks && (
              <div className="add-expense-link">
                <Link data-cy="expense-link" to={`/moves/${moveId}/ppm-expenses`}>
                  <FontAwesomeIcon className="icon link-blue" icon={faPlusCircle} /> Add expense
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
    weightTicketDocs: selectPPMCloseoutDocumentsForMove(state, moveId, ['WEIGHT_TICKET_SET']),
  };
}

const mapDispatchToProps = {
  selectPPMCloseoutDocumentsForMove,
  getMoveDocumentsForMove,
};

export default connect(mapStateToProps, mapDispatchToProps)(DocumentsUploaded);
