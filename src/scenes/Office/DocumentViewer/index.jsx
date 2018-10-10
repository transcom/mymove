import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get } from 'lodash';
import qs from 'query-string';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import Alert from 'shared/Alert';
import { PanelField } from 'shared/EditablePanel';
import { loadMoveDependencies } from '../ducks.js';
import { Tab, Tabs, TabList, TabPanel } from 'react-tabs';
import PrivateRoute from 'shared/User/PrivateRoute';
import { Switch, Redirect, Link } from 'react-router-dom';

import DocumentUploadViewer from 'shared/DocumentViewer/DocumentUploadViewer';
import DocumentList from 'shared/DocumentViewer/DocumentList';
import DocumentDetailPanel from './DocumentDetailPanel';

import DocumentUploader from './DocumentUploader';
import {
  selectAllDocumentsForMove,
  getMoveDocumentsForMove,
} from 'shared/Entities/modules/moveDocuments';
import { stringifyName } from 'shared/utils/serviceMember';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPlusCircle from '@fortawesome/fontawesome-free-solid/faPlusCircle';

import './index.css';
class DocumentViewer extends Component {
  componentDidMount() {
    //this is probably overkill, but works for now
    this.props.loadMoveDependencies(this.props.match.params.moveId);
    this.props.getMoveDocumentsForMove(this.props.match.params.moveId);
  }
  componentWillUpdate() {
    document.title = 'Document Viewer';
  }
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

    // Parse query string parameters
    const queryValues = qs.parse(this.props.location.search);

    const defaultTabIndex =
      this.props.match.params.moveDocumentId !== 'new' ? 1 : 0;
    if (
      !this.props.loadDependenciesHasSuccess &&
      !this.props.loadDependenciesHasError
    )
      return <LoadingPlaceholder />;
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
              <PrivateRoute
                exact
                path={defaultPath}
                render={() => <Redirect replace to={newUrl} />}
              />
              <PrivateRoute
                path={newPath}
                moveId={move.id}
                render={() => {
                  return (
                    <DocumentUploader
                      moveId={move.id}
                      moveDocumentType={queryValues.move_document_type}
                      location={this.props.location}
                    />
                  );
                }}
              />
              <PrivateRoute
                path={documentPath}
                component={DocumentUploadViewer}
              />
              <PrivateRoute
                path={defaultUrl}
                render={() => <div> document viewer coming soon</div>}
              />
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
                <Tab className="title nav-tab">
                  All Documents ({numMoveDocs})
                </Tab>
                {/* TODO: Handle routing of /new route better */}
                {this.props.match.params.moveDocumentId &&
                  this.props.match.params.moveDocumentId !== 'new' && (
                    <Tab className="title nav-tab">Details</Tab>
                  )}
              </TabList>

              <TabPanel>
                <div className="pad-ns">
                  <span className="status">
                    <FontAwesomeIcon
                      className="icon link-blue"
                      icon={faPlusCircle}
                    />
                  </span>
                  <Link to={newUrl}>Upload new document</Link>
                </div>
                <div>
                  {' '}
                  <DocumentList
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

DocumentViewer.propTypes = {
  loadMoveDependencies: PropTypes.func.isRequired,
};

const mapStateToProps = state => ({
  swaggerError: state.swaggerInternal.hasErrored,
  orders: state.office.officeOrders || {},
  move: get(state, 'office.officeMove', {}),
  moveDocuments: selectAllDocumentsForMove(
    state,
    get(state, 'office.officeMove.id', ''),
  ),
  serviceMember: state.office.officeServiceMember || {},
  loadDependenciesHasSuccess: state.office.loadDependenciesHasSuccess,
  loadDependenciesHasError: state.office.loadDependenciesHasError,
  error: state.office.error,
});

const mapDispatchToProps = dispatch =>
  bindActionCreators(
    { loadMoveDependencies, getMoveDocumentsForMove },
    dispatch,
  );

export default connect(mapStateToProps, mapDispatchToProps)(DocumentViewer);
