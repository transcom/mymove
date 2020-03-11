import { shallow } from 'enzyme';
import { DocumentDetailDisplay } from './DocumentDetailPanel';
import React from 'react';
import { WEIGHT_TICKET_SET_TYPE } from '../../../shared/constants';

describe('DocumentDetailDisplay', () => {
  const renderDocumentDetailDisplay = ({
    isExpenseDocument = true,
    isWeightTicketDocument = true,
    moveDocument = {},
    moveDocSchema = {
      properties: { title: { enum: false }, move_document_type: false },
      required: [],
      type: 'string type',
    },
    isStorageExpenseDocument = false,
  }) =>
    shallow(
      <DocumentDetailDisplay
        isExpenseDocument={isExpenseDocument}
        isWeightTicketDocument={isWeightTicketDocument}
        moveDocument={moveDocument}
        moveDocSchema={moveDocSchema}
        isStorageExpenseDocument={isStorageExpenseDocument}
      />,
    );

  describe('weight ticket document display view', () => {
    it('includes common document info', () => {
      const moveDocument = {
        id: 'id',
        move_id: 'id',
        title: 'My Title',
        move_document_type: 'WEIGHT_TICKET_SET',
        notes: 'This is a note',
        status: 'AWAITING_REVIEW',
        document: {
          id: 'an id',
          move_document_id: 'another id',
          service_member_id: 'another id2',
          uploads: [
            {
              id: 'id',
              url: 'url',
              filename: 'file here',
              content_type: 'json',
              created_at: '2018-09-27 08:14:38.702434',
            },
          ],
        },
      };

      const documentDisplay = renderDocumentDetailDisplay({ moveDocument });
      console.log(
        documentDisplay
          .find('[data-cy="title"]')
          .dive()
          .dive()
          .debug(),
      );
      expect(documentDisplay.find('[data-cy="panel-subhead"]').text()).toContain(moveDocument.title);
      expect(documentDisplay.find('[data-cy="uploaded-at"]').text()).toEqual('Uploaded 27-Sep-18');
      expect(
        documentDisplay
          .find('[data-cy="title"]')
          .dive()
          .dive()
          .find('SwaggerValue')
          .dive()
          .text(),
      ).toEqual(moveDocument.title);
      expect(
        documentDisplay
          .find('[data-cy="move-document-type"]')
          .dive()
          .dive()
          .find('SwaggerValue')
          .dive()
          .text(),
      ).toEqual(moveDocument.move_document_type);
      expect(documentDisplay.find('[data-cy="status"]').text()).toEqual(moveDocument.notes);
      expect(documentDisplay.find('[data-cy="notes"]').text()).toEqual(moveDocument.notes);
    });
    it('includes weight ticket-specific fields', () => {
      const moveDocument = {
        emptyWeight: 2200,
        fullWeight: 3500,
      };
      const documentDisplay = renderDocumentDetailDisplay({ moveDocument });
      expect(documentDisplay.find('[data-cy="empty-weight"]').text()).toEqual(moveDocument.emptyWeight);
      expect(documentDisplay.find('[data-cy="full-weight"]').text()).toEqual(moveDocument.fullWeight);
    });

    describe('is car or car and trailer', () => {
      it('includes the make and model fields ', () => {
        const moveDocument = {
          vehicleMake: 'Honda',
          vehicleModel: 'Civic',
          weightTicketSetType: WEIGHT_TICKET_SET_TYPE.CAR,
        };
        const documentDisplay = renderDocumentDetailDisplay({ moveDocument });
        expect(documentDisplay.find('[data-cy="vehicle-make"]').text()).toEqual(moveDocument.vehicleMake);
        expect(documentDisplay.find('[data-cy="vehicle-model"]').text()).toEqual(moveDocument.vehicleModel);
      });
    });

    describe('a box truck type weight ticket', () => {
      it('includes vehicle nickname', () => {
        const moveDocument = {
          vehicleNickname: 'Civic',
          weightTicketSetType: WEIGHT_TICKET_SET_TYPE.BOX_TRUCK,
        };
        const documentDisplay = renderDocumentDetailDisplay({ moveDocument });
        expect(documentDisplay.find('[data-cy="vehicle-nickname"]').text()).toEqual(moveDocument.vehicleNickname);
      });
    });
  });
});
