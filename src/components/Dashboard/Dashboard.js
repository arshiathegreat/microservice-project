// Dashboard.js
import React, { useState, useEffect } from 'react';
import { Container, Row, Col, Button, Table } from 'react-bootstrap';


export default function Dashboard({ token, handleLogout }) {

  const [data, setData] = useState([]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch('http://localhost:8081/api/s1/users', {
          method: 'GET',
          headers: {
            'Authorization': `${token}`
          }
        });
        if (response.ok) {
          const result = await response.json();
          setData(result.data);
        } else {
          console.error('Failed to fetch data');
        }
      } catch (error) {
        console.error('Error:', error);
      }
    };

    fetchData();
  }, [token]);

  return (
    <Container>
      <Row className="mt-3">
        <Col xs={4}>
          <h4>Welcome Admin</h4>
        </Col>
        <Col xs={8} className="text-right">
          <Button variant="danger" onClick={handleLogout}>Logout</Button>
        </Col>
      </Row>
      <Row className="mt-3">
      <Col>
        <Table striped bordered hover>
          <thead>
            <tr>
              
              <th>Name</th>
              <th>Email</th>
              <th>Role</th>
            </tr>
          </thead>
          <tbody>
            {data.map(item => (
              <tr key={item.id}>
                <td>{item.username}</td>
                <td>{item.email}</td>
                <td>{item.role}</td>
              </tr>
            ))}
          </tbody>
        </Table>
      </Col>
    </Row>
    </Container>
  );
}
