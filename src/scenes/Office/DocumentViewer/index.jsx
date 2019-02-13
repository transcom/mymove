import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { includes, get } from 'lodash';
import qs from 'query-string';

import { createMoveDocument } from 'shared/Entities/modules/moveDocuments';
import { createMovingExpenseDocument } from 'shared/Entities/modules/movingExpenseDocuments';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import Alert from 'shared/Alert';
import { PanelField } from 'shared/EditablePanel';
import { loadMoveDependencies } from '../ducks.js';
import { selectServiceMemberForOrders } from 'shared/Entities/modules/serviceMembers';
import { selectOrdersForMove } from 'shared/Entities/modules/orders';
import { Tab, Tabs, TabList, TabPanel } from 'react-tabs';
import PrivateRoute from 'shared/User/PrivateRoute';
import { Switch, Redirect, Link } from 'react-router-dom';

import DocumentUploadViewer from 'shared/DocumentViewer/DocumentUploadViewer';
import DocumentList from 'shared/DocumentViewer/DocumentList';
import { selectPPMForMove } from 'shared/Entities/modules/ppms';

import DocumentUploader from 'shared/DocumentViewer/DocumentUploader';
import { selectAllDocumentsForMove, getMoveDocumentsForMove } from 'shared/Entities/modules/moveDocuments';
import { stringifyName } from 'shared/utils/serviceMember';
import { convertDollarsToCents } from 'shared/utils';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPlusCircle from '@fortawesome/fontawesome-free-solid/faPlusCircle';

import DocumentDetailPanel from './DocumentDetailPanel';

import './index.css';
class DocumentViewer extends Component {
  static propTypes = {
    docTypes: PropTypes.arrayOf(PropTypes.string).isRequired,
    loadMoveDependencies: PropTypes.func.isRequired,
    getMoveDocumentsForMove: PropTypes.func.isRequired,
    genericMoveDocSchema: PropTypes.object.isRequired,
    moveDocSchema: PropTypes.object.isRequired,
    moveDocuments: PropTypes.arrayOf(PropTypes.object),
    location: PropTypes.object.isRequired,
  };
  componentDidMount() {
    //this is probably overkill, but works for now
    this.props.loadMoveDependencies(this.props.match.params.moveId);
    this.props.getMoveDocumentsForMove(this.props.match.params.moveId);
  }
  componentWillUpdate() {
    document.title = 'Document Viewer';
  }

  get getDocumentUploaderProps() {
    const { docTypes, location, genericMoveDocSchema, moveDocSchema } = this.props;
    // Parse query string parameters
    const moveDocumentType = qs.parse(location.search).moveDocumentType;

    const initialValues = {};
    // Verify the provided doc type against the schema
    if (includes(docTypes, moveDocumentType)) {
      initialValues.move_document_type = moveDocumentType;
    }

    return {
      form: 'move_document_upload',
      isPublic: false,
      onSubmit: this.handleSubmit,
      genericMoveDocSchema,
      initialValues,
      location,
      moveDocSchema,
    };
  }

  handleSubmit = (uploadIds, formValues) => {
    const {
      currentPpm,
      move: { id: moveId },
    } = this.props;
    const {
      title,
      moving_expense_type: movingExpenseType,
      move_document_type: moveDocumentType,
      requested_amount_cents: requestedAmountCents,
      payment_method: paymentMethod,
      notes,
    } = formValues;
    const personallyProcuredMoveId = currentPpm ? currentPpm.id : null;
    if (get(formValues, 'move_document_type', false) === 'EXPENSE') {
      return this.props.createMovingExpenseDocument({
        moveId,
        personallyProcuredMoveId,
        uploadIds,
        title,
        movingExpenseType,
        moveDocumentType,
        requestedAmountCents: convertDollarsToCents(requestedAmountCents),
        paymentMethod,
        notes,
      });
    }
    return this.props.createMoveDocument({
      moveId,
      personallyProcuredMoveId,
      uploadIds,
      title,
      moveDocumentType,
      notes,
    });
  };
  render() {
    const { serviceMember, move, moveDocuments } = this.props;
    const numMoveDocs = moveDocuments ? moveDocuments.length : 0;
    const name = stringifyName(serviceMember);

    // urls: has full url with IDs
    const defaultUrl = move ? `/moves/${move.id}/documents` : '';
    const newUrl = move ? `/moves/${move.id}/documents/new` : '';

    // paths: has placeholders (e.g. ":moveId")
    const defaultPath = `/moves/:moveId/documents`;
    const newPath = `/moves/:moveId/documents/new`;
    const documentPath = `/moves/:moveId/documents/:moveDocumentId`;
    const currentMoveDocumentId = this.props.match.params.moveDocumentId;

    const defaultTabIndex = this.props.match.params.moveDocumentId !== 'new' ? 1 : 0;
    if (!this.props.loadDependenciesHasSuccess && !this.props.loadDependenciesHasError) return <LoadingPlaceholder />;
    if (this.props.loadDependenciesHasError)
      return (
        <div className="usa-grid">
          <div className="usa-width-one-whole error-message">
            <Alert type="error" heading="An error occurred">
              Something went wrong contacting the server.
            </Alert>
          </div>
        </div>
      );
    return (
      <div className="usa-grid doc-viewer">
        <div className="usa-width-two-thirds">
          <div className="tab-content">
            <Switch>
              <PrivateRoute exact path={defaultPath} render={() => <Redirect replace to={newUrl} />} />
              <PrivateRoute
                path={newPath}
                moveId={move.id}
                render={() => {
                  return <DocumentUploader {...this.getDocumentUploaderProps} />;
                }}
              />
              <PrivateRoute path={documentPath} component={DocumentUploadViewer} />
              <PrivateRoute path={defaultUrl} render={() => <div> document viewer coming soon</div>} />
            </Switch>
          </div>
        </div>
        <div className="usa-width-one-third">
          <h3>{name}</h3>
          <PanelField title="Move Locator">{move.locator}</PanelField>
          <PanelField title="DoD ID">{serviceMember.edipi}</PanelField>
          <div className="tab-content">
            <Tabs defaultIndex={defaultTabIndex}>
              <TabList className="doc-viewer-tabs">
                <Tab className="title nav-tab">All Documents ({numMoveDocs})</Tab>
                {/* TODO: Handle routing of /new route better */}
                {this.props.match.params.moveDocumentId &&
                  this.props.match.params.moveDocumentId !== 'new' && <Tab className="title nav-tab">Details</Tab>}
              </TabList>

              <TabPanel>
                <div className="pad-ns">
                  <span className="status">
                    <FontAwesomeIcon className="icon link-blue" icon={faPlusCircle} />
                  </span>
                  <Link to={newUrl}>Upload new document</Link>
                </div>
                <div>
                  {' '}
                  <DocumentList
                    currentMoveDocumentId={currentMoveDocumentId}
                    detailUrlPrefix={`/moves/${move.id}/documents`}
                    moveDocuments={moveDocuments}
                  />
                </div>
              </TabPanel>

              {this.props.match.params.moveDocumentId &&
                this.props.match.params.moveDocumentId !== 'new' && (
                  <TabPanel>
                    <DocumentDetailPanel
                      className="document-viewer"
                      moveDocumentId={this.props.match.params.moveDocumentId}
                      moveId={move.id}
                      title=""
                    />
                  </TabPanel>
                )}
            </Tabs>
          </div>
        </div>
      </div>
    );
  }
}

const mapStateToProps = state => {
  const move = get(state, 'office.officeMove') || {};
  const serviceMember = selectServiceMemberForOrders(state, move.orders_id);

  return {
    genericMoveDocSchema: get(state, 'swaggerInternal.spec.definitions.CreateGenericMoveDocumentPayload', {}),
    moveDocSchema: get(state, 'swaggerInternal.spec.definitions.MoveDocumentPayload', {}),
    currentPpm: selectPPMForMove(state, move.id),
    docTypes: get(state, 'swaggerInternal.spec.definitions.MoveDocumentType.enum', []),
    orders: selectOrdersForMove(state, move.id),
    move,
    moveDocuments: selectAllDocumentsForMove(state, get(state, 'office.officeMove.id', '')),
    serviceMember,
    loadDependenciesHasSuccess: state.office.loadDependenciesHasSuccess,
    loadDependenciesHasError: state.office.loadDependenciesHasError,
    error: state.office.error,
  };
};

const mapDispatchToProps = {
  createMoveDocument,
  createMovingExpenseDocument,
  loadMoveDependencies,
  getMoveDocumentsForMove,
};

export default connect(mapStateToProps, mapDispatchToProps)(DocumentViewer);
