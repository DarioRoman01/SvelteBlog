<script lang="ts">
  import { MaterialApp, Icon, Container, Row, Col, Button } from "svelte-materialify";
  import { mdiAccountCircle } from "@mdi/js";
  import type { Profile } from "../../requests/profile";

  let profile: Profile;
</script>

<MaterialApp>
  <Container>
    <Row>
      <Col>
        <div>
          <Icon path={mdiAccountCircle}/>
          <Button class="pink lighten-3">
            follow
          </Button>
        </div>
      </Col>
      <Col>
        <div>
          <p>{profile.username}</p>
          <p>{profile.followers}</p>
          <p>{profile.posted}</p>
        </div>
      </Col>
    </Row>
  </Container>
</MaterialApp>