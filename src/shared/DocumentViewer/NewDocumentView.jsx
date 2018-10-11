import React, { Component } from 'react';
import { func, object, array, string, shape } from 'prop-types';
import { Link } from 'react-router-dom';

import { PanelField } from 'shared/EditablePanel';
import { Tab, Tabs, TabList, TabPanel } from 'react-tabs';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPlusCircle from '@fortawesome/fontawesome-free-solid/faPlusCircle';
import DocumentUploader from 'shared/DocumentViewer/DocumentUploader';

import './index.css';

import DocumentList from './DocumentList';

class NewDocumentView extends Component {
  static propTypes = {
    createShipmentDocument: func.isRequired,
    genericMoveDocSchema: object.isRequired,
    moveDocSchema: object.isRequired,
    moveDocuments: array.isRequired,
    moveLocator: string.isRequired,
    onDidMount: func.isRequired,
    shipmentId: string.isRequired,
    serviceMember: shape({
      edipi: string.isRequired,
      name: string.isRequired,
    }),
  };
  componentDidMount() {
    const { onDidMount } = this.props;
    onDidMount();
  }

  handleSubmit = (upload_ids, formValues) => {
    const { createShipmentDocument, shipmentId } = this.props;
    const { move_document_type, title, notes } = formValues;
    const createGenericMoveDocument = {
      shipmentId,
      upload_ids,
      move_document_type,
      title,
      notes,
    };

    return createShipmentDocument(shipmentId, createGenericMoveDocument);
  };

  render() {
    const {
      genericMoveDocSchema,
      moveDocSchema,
      moveDocuments,
      moveLocator,
      shipmentId,
      serviceMember: { edipi, name },
    } = this.props;
    const newDocumentUrl = `/shipments/${shipmentId}/documents/new`;
    const documentDetailUrlPrefix = `/shipments/${shipmentId}/documents`;

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
                  <Link to={newDocumentUrl}>Upload new document</Link>
                </div>
                <div>
                  {' '}
                  <DocumentList
                    detailUrlPrefix={documentDetailUrlPrefix}
                    moveDocuments={moveDocuments}
                  />
                </div>
              </TabPanel>
            </Tabs>
          </div>
        </div>
      </div>
    );
  }
}

export default NewDocumentView;
