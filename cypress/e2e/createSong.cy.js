describe('Create Song API Test', () => {
    it('should create a new song successfully', () => {
      // Prepare song data
      const songData = {
        title: "Imagine",
        artist: "John Lennon",
        genre: "Rock"
      };
  
      // Send a POST request to the /songs endpoint
      cy.request('POST', '/songs', songData)
        .then((response) => {
          // Assert that the response status is 200
          expect(response.status).to.eq(200);
  
          // Assert the response body contains the created song
          expect(response.body).to.have.property('title', songData.title);
          expect(response.body).to.have.property('artist', songData.artist);
          expect(response.body).to.have.property('genre', songData.genre);
        });
    });
  
    it('should return 400 for missing artist field', () => {
      const invalidData = { title: "No Artist" };
  
      cy.request({
        method: 'POST',
        url: '/songs',
        failOnStatusCode: false, // Don't fail the test if the response status isn't 2xx
        body: invalidData,
      }).then((response) => {
        // Assert that the response status is 400
        expect(response.status).to.eq(400);
  
        // Assert the error message
        expect(response.body).to.include('Missing required fields');
      });
    });

    it('should return 400 for missing title field', () => {
        const invalidData = {title: "No title"};

        cy.request({
            method: 'POST',
            url: '/songs',
            failOnStatusCode: false,
            body: invalidData,
        }).then((response) => {
            expect(response.status).to.eq(400);
            expect(response.body).to.include('Missing required fields')
        });
    });

    it('should return 400 for missing genre field', () => {
        const invalidData = {title: "No genre"};

        cy.request({
            method: 'POST',
            url: '/songs',
            failOnStatusCode: false,
            body: invalidData,
        }).then((response) => {
            expect(response.status).to.eq(400);
            expect(response.body).to.include('Missing required fields')
        });
    });
  });
  