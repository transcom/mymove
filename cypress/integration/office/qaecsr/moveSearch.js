import { TOOOfficeUserType } from '../../../support/constants';

describe('QAE/CSR Move Search', () => {
  before(() => {
    cy.prepareOfficeApp();
  });

  beforeEach(() => {});

  it('is able to search by move code', () => {
    cy.visit('/qaecsr/search');
  });
});
