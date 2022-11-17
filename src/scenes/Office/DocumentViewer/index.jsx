import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { includes, get, isEmpty } from 'lodash';
import qs from 'query-string';

import { selectMove } from 'shared/Entities/modules/moves';
import { createMovingExpenseDocument } from 'shared/Entities/modules/movingExpenseDocuments';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import Alert from 'shared/Alert';
import { PanelField } from 'shared/EditablePanel';
import { loadMove, loadMoveLabel } from 'shared/Entities/modules/moves';
import { getRequestStatus } from 'shared/Swagger/selectors';
import { loadServiceMember, selectServiceMember } from 'shared/Entities/modules/serviceMembers';
import { Tab, Tabs, TabList, TabPanel } from 'react-tabs';
import PrivateRoute from 'containers/PrivateRoute';
import { Switch, Redirect } from 'react-router-dom-old';

import DocumentUploadViewer from 'shared/DocumentViewer/DocumentUploadViewer';
import DocumentList from 'shared/DocumentViewer/DocumentList';
import { selectActivePPMForMove } from 'shared/Entities/modules/ppms';

import DocumentUploader from 'shared/DocumentViewer/DocumentUploader';
import {
  selectAllDocumentsForMove,
  getMoveDocumentsForMove,
  createMoveDocument,
} from 'shared/Entities/modules/moveDocuments';
import { stringifyName } from 'shared/utils/serviceMember';
import { convertDollarsToCents } from 'shared/utils';
import { roleTypes } from 'constants/userRoles';

import DocumentDetailPanel from './DocumentDetailPanel';

import './index.css';

class DocumentViewer extends Component {
  componentDidMount() {
    const { moveId } = this.props;
    this.props.loadMove(moveId);
    this.props.getMoveDocumentsForMove(moveId);
  }

  componentDidUpdate(prevProps) {
    const { serviceMemberId } = this.props;
    if (serviceMemberId !== prevProps.serviceMemberId) {
      this.props.loadServiceMember(serviceMemberId);
    }
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
    const { currentPpm, moveId } = this.props;
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
        missingReceipt: false,
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
    const { serviceMember, moveId, moveDocumentId, moveDocuments, moveLocator } = this.props;
    const numMoveDocs = moveDocuments ? moveDocuments.length : 0;
    const name = stringifyName(serviceMember);
    document.title = `Document Viewer for ${name}`;

    // urls: has full url with IDs
    const defaultUrl = `/moves/${moveId}/documents`;
    const newUrl = `/moves/${moveId}/documents/new`;

    // paths: has placeholders (e.g. ":moveId")
    const defaultPath = `/moves/:moveId/documents`;
    const newPath = `/moves/:moveId/documents/new`;
    const documentPath = `/moves/:moveId/documents/:moveDocumentId`;

    const defaultTabIndex = this.props.match.params.moveDocumentId !== 'new' ? 1 : 0;
    if (!this.props.loadDependenciesHasSuccess && !this.props.loadDependenciesHasError) return <LoadingPlaceholder />;
    if (this.props.loadDependenciesHasError)
      return (
        <div className="grid-container-widescreen usa-prose">
          <div className="grid-row">
            <div className="grid-col-12 error-message">
              <Alert type="error" heading="An error occurred">
                Something went wrong contacting the server.
              </Alert>
            </div>
          </div>
        </div>
      );
    return (
      <div className="grid-container-widescreen usa-prose">
        <div className="grid-row grid-gap doc-viewer">
          <div className="grid-col-8">
            <div className="tab-content">
              <Switch>
                <PrivateRoute
                  exact
                  path={defaultPath}
                  render={() => <Redirect replace to={newUrl} />}
                  requiredRoles={[roleTypes.PPM]}
                />
                <PrivateRoute
                  path={newPath}
                  moveId={moveId}
                  render={() => {
                    return <DocumentUploader {...this.getDocumentUploaderProps} />;
                  }}
                  requiredRoles={[roleTypes.PPM]}
                />
                <PrivateRoute path={documentPath} component={DocumentUploadViewer} requiredRoles={[roleTypes.PPM]} />
                <PrivateRoute
                  path={defaultUrl}
                  render={() => <div> document viewer coming soon</div>}
                  requiredRoles={[roleTypes.PPM]}
                />
              </Switch>
            </div>
          </div>
          <div className="grid-col-4">
            <h3>{name}</h3>
            <PanelField title="Move Locator">{moveLocator}</PanelField>
            <PanelField title="DoD ID">{serviceMember.edipi}</PanelField>
            <div className="tab-content">
              <Tabs defaultIndex={defaultTabIndex}>
                <TabList className="doc-viewer-tabs">
                  <Tab className="title nav-tab">All Documents ({numMoveDocs})</Tab>
                  {/* TODO: Handle routing of /new route better */}
                  {moveDocumentId && moveDocumentId !== 'new' && <Tab className="title nav-tab">Details</Tab>}
                </TabList>

                <TabPanel>
                  <div>
                    {' '}
                    <DocumentList
                      currentMoveDocumentId={moveDocumentId}
                      detailUrlPrefix={`/moves/${moveId}/documents`}
                      moveDocuments={moveDocuments}
                      uploadDocumentUrl={newUrl}
                      moveId={moveId}
                    />
                  </div>
                </TabPanel>

                {!isEmpty(moveDocuments) && moveDocumentId && moveDocumentId !== 'new' && (
                  <TabPanel>
                    <DocumentDetailPanel
                      className="document-viewer"
                      moveDocumentId={moveDocumentId}
                      moveId={moveId}
                      title=""
                    />
                  </TabPanel>
                )}
              </Tabs>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

DocumentViewer.propTypes = {
  docTypes: PropTypes.arrayOf(PropTypes.string).isRequired,
  loadMove: PropTypes.func.isRequired,
  getMoveDocumentsForMove: PropTypes.func.isRequired,
  genericMoveDocSchema: PropTypes.object.isRequired,
  moveDocSchema: PropTypes.object.isRequired,
  moveDocuments: PropTypes.arrayOf(PropTypes.object),
  location: PropTypes.object.isRequired,
};

const mapStateToProps = (state, ownProps) => {
  const { moveId, moveDocumentId } = ownProps.match.params;
  const move = selectMove(state, moveId);
  const moveLocator = move.locator;
  const serviceMemberId = move.service_member_id;
  const serviceMember = selectServiceMember(state, serviceMemberId);
  const loadMoveRequest = getRequestStatus(state, loadMoveLabel);

  return {
    genericMoveDocSchema: get(state, 'swaggerInternal.spec.definitions.CreateGenericMoveDocumentPayload', {}),
    moveDocSchema: get(state, 'swaggerInternal.spec.definitions.MoveDocumentPayload', {}),
    currentPpm: selectActivePPMForMove(state, moveId),
    docTypes: get(state, 'swaggerInternal.spec.definitions.MoveDocumentType.enum', []),
    moveId,
    moveLocator,
    moveDocumentId,
    moveDocuments: selectAllDocumentsForMove(state, moveId),
    serviceMember,
    serviceMemberId,
    loadDependenciesHasSuccess: loadMoveRequest.isSuccess,
    loadDependenciesHasError: loadMoveRequest.error,
  };
};

const mapDispatchToProps = {
  createMoveDocument,
  createMovingExpenseDocument,
  loadMove,
  loadServiceMember,
  getMoveDocumentsForMove,
};

export default connect(mapStateToProps, mapDispatchToProps)(DocumentViewer);
