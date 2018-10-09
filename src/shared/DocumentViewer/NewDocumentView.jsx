import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router-dom';

import { PanelField } from 'shared/EditablePanel';
import { Tab, Tabs, TabList, TabPanel } from 'react-tabs';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPlusCircle from '@fortawesome/fontawesome-free-solid/faPlusCircle';
import DocumentUploader from 'shared/DocumentViewer/DocumentUploader';
import DocumentDetailPanelView from 'shared/DocumentViewer/DocumentDetailPanelView';

import { createShipmentDocumentLabel } from 'shared/Entities/modules/uploads';
import './index.css';

import DocumentList from './DocumentList';
import DocumentContent from './DocumentContent';

class NewDocumentView extends Component {
  static propTypes = {
    documentDetailUrlPrefix: PropTypes.string.isRequired,
    shipmentId: PropTypes.string.isRequired,
    createShipmentDocument: PropTypes.func.isRequired,
  };
  componentDidMount() {
    const { onDidMount } = this.props;
    onDidMount();
  }

  handleSubmit = (uploadIds, formValues) => {
    const { shipmentId } = this.props;
    const { move_document_type, title, notes } = formValues;

    const createGenericMoveDocument = {
      shipmentId,
      uploadIds,
      move_document_type,
      title,
      notes,
    };

    return this.props.createShipmentDocument(
      createShipmentDocumentLabel,
      shipmentId,
      createGenericMoveDocument,
    );
  };

  render() {
    const {
      documentDetailUrlPrefix,
      genericMoveDocSchema,
      moveDocSchema,
      moveDocuments,
      moveLocator,
      // newDocumentUrl,
      serviceMember: { edipi, name },
    } = this.props;
    return (
      <div className="usa-grid doc-viewer">
        <div className="usa-width-two-thirds">
          <div className="tab-content">
            <div className="document-contents">
              <DocumentUploader
                form="shipmment-documents"
                initialValues={{}}
                genericMoveDocSchema={genericMoveDocSchema}
                moveDocSchema={moveDocSchema}
                onSubmit={this.handleSubmit}
                isPublic={true}
              />
            </div>
          </div>
        </div>
        <div className="usa-width-one-third">
          <h3>{name}</h3>
          <PanelField title="Move Locator">{moveLocator}</PanelField>
          <PanelField title="DoD ID">{edipi}</PanelField>
          <div className="tab-content">
            <Tabs defaultIndex={0}>
              <TabList className="doc-viewer-tabs">
                <Tab className="title nav-tab">
                  All Documents ({moveDocuments.length})
                </Tab>
              </TabList>

              <TabPanel>
                <div className="pad-ns">
                  <span className="status">
                    <FontAwesomeIcon
                      className="icon link-blue"
                      icon={faPlusCircle}
                    />
                  </span>
                  {/* TODO Do we want to have the new doc URL on the new page? */}
                  {/* <Link to={newDocumentUrl}>Upload new document</Link> */}
                </div>
              </TabPanel>
              <TabPanel>
                <div className="pad-ns">
                  <span className="status">
                    <FontAwesomeIcon
                      className="icon link-blue"
                      icon={faPlusCircle}
                    />
                  </span>
                  {/* <Link to={newDocumentUrl}>Upload new document</Link> */}
                </div>
                <div>
                  {' '}
                  {/* <DocumentList
                    detailUrlPrefix={documentDetailUrlPrefix}
                    moveDocuments={moveDocuments}
                  /> */}
                </div>
              </TabPanel>
            </Tabs>
          </div>
        </div>
      </div>
    );
  }
}

NewDocumentView.propTypes = {
  // documentDetailUrlPrefix: PropTypes.string.isRequired,
  // moveDocumentSchema: PropTypes.object.isRequired,
  moveDocuments: PropTypes.array.isRequired,
  moveLocator: PropTypes.string.isRequired,
  // newDocumentUrl: PropTypes.string.isRequired,
  onDidMount: PropTypes.func.isRequired,
  serviceMember: PropTypes.shape({
    edipi: PropTypes.string.isRequired,
    name: PropTypes.string.isRequired,
  }),
};

export default NewDocumentView;
